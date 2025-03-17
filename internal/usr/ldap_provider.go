package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"
	"stoke/internal/ent/grouplink"
	"stoke/internal/ent/predicate"
	"stoke/internal/tel"
	"strings"
	"text/template"

	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

var (
	LDAPNotFoundError = errors.New("No results")
	LDAPError = errors.New("Error communicating with LDAP")
)

type LDAPConnector interface {
	Connect(string, ...ldap.DialOpt) (ldap.Client, error)
}

type ldapUserProvider struct {
	Name                  string
	ServerURL             string
	BindUserDN            string
	BindUserPassword      string

	GroupSearchRoot       string
	GroupFilter           *template.Template
	GroupAttribute        string

	UserSearchRoot        string
	UserFilter            *template.Template

	FirstNameField        string
	LastNameField         string
	EmailField            string

	SearchTimeout         int

	DialOpts 							[]ldap.DialOpt

	connector             LDAPConnector
}

type templateValues struct {
	Username string
	UserDN   string
}

type ldapDialer struct {}

// Creates a NewLDAPUserProvider
func NewLDAPUserProvider(name, url, bindDN, bindPass, groupSearch, groupAttribute, userSearch, fnameField, lnameField, emailField string, searchTimeout int, groupFilter, userFilter *template.Template, dialOpts ...ldap.DialOpt) *ldapUserProvider {
	return &ldapUserProvider{
		Name:             name,
		ServerURL:        url,
		BindUserDN:       bindDN,
		BindUserPassword: bindPass,
		GroupSearchRoot:  groupSearch,
		GroupFilter:      groupFilter,
		GroupAttribute:   groupAttribute,
		UserSearchRoot:   userSearch,
		UserFilter:       userFilter,
		FirstNameField:   fnameField,
		LastNameField:    lnameField,
		EmailField:       emailField,
		SearchTimeout:    searchTimeout,
		DialOpts:         dialOpts,
		connector:        ldapDialer{},
	}
}

// Connect using the ldap.DialURL function
func (ldapDialer) Connect(url string, dialOpts ...ldap.DialOpt) (ldap.Client, error) {
	return ldap.DialURL(url, dialOpts...)
}

// Set the ldap connector to use.
// Should only be needed for testing, but could also be used to alter connection behaviour
func (l *ldapUserProvider) SetConnector(c LDAPConnector) {
	l.connector = c
}

// GetUserClaims looks up claims that are associated in LDAP
func (l *ldapUserProvider) UpdateUserClaims(username, password string, ctx context.Context) (*ent.User, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("url", l.ServerURL).
		Str("username", username).
		Str("bindUserDN", l.BindUserDN).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "ldapUserProvider.UpdateUserClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Msg("Getting user claims")

	conn, err := l.connector.Connect(l.ServerURL, l.DialOpts...)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not connect to LDAP server")
		return nil, AuthSourceError
	}
	defer conn.Close()

	if err:= conn.Bind(l.BindUserDN, l.BindUserPassword); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Bind user authentication failed")
		return nil, AuthSourceError
	}

	usr, groupLinks, err := l.getOrCreateUser(username, password, conn, ctx)
	if errors.Is(LDAPNotFoundError, err) || errors.Is(LDAPError, err) {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not find User in ldap")
		return nil, UserNotFoundError

	} else if errors.Is(NoLinkedGroupsError, err) {
		if usr == nil {
			return nil, AuthenticationError
		}
	} else if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not get or create user")
		return nil, err
	}

	// Update local database with claims that the user has
	// LDAP user groups that the user already has
	addClaimGroups, delClaimGroups := findGroupChanges(usr, groupLinks)

	if usr, err = applyGroupChanges(addClaimGroups, delClaimGroups, usr, ctx); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Failed to add LDAP groups to local user")
		return nil, err
	}

	usr, err = retreiveLocalUser(usr.Username, ctx)
	if len(usr.Edges.ClaimGroups) == 0 {
		return nil, NoLinkedGroupsError
	} 
	return usr, err
}

// Creates the user if it exists in LDAP
func (l *ldapUserProvider) getOrCreateUser(username, password string, conn ldap.Client, ctx context.Context) (*ent.User, ent.GroupLinks, error) {
	logger := zerolog.Ctx(ctx).With().
			Str("url", l.ServerURL).
			Str("user", username).
			Logger()

	userEntry, err := l.getLDAPUser(username, password, conn, ctx)
	if err != nil{
		return nil, nil, err
	}

	usr, dbErr := retreiveLocalUser(username, ctx)
	if dbErr != nil && !ent.IsNotFound(dbErr) {
		return nil, nil, dbErr
	}

	if ent.IsNotFound(dbErr) {
		usr, err = l.createLocalUser(username, userEntry, ctx)
		if err != nil {
			logger.Error().
				Err(err).
				Str("userDN", userEntry.DN).
				Msg("Could not create local user")
			return nil, nil, err
		}
	}

	groupLinks, err := l.getUserLDAPGroupLinks(username, userEntry.DN, conn, ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Str("userDN", userEntry.DN).
			Msg("Could not get users linked LDAP groups")
		return nil, nil, err
	}

	if len(groupLinks) == 0 {
		logger.Warn().
			Err(err).
			Str("userDN", userEntry.DN).
			Msg("User has no linked groups")
		return usr, nil, NoLinkedGroupsError
	}

	logger.Debug().
		Int("numGroups", len(groupLinks)).
		Msg("Found user with group links")

	return usr, groupLinks, nil
}

func (l *ldapUserProvider) getLDAPUser(username, password string, conn ldap.Client, ctx context.Context) (*ldap.Entry, error) {
	result, err := l.ldapSearch(
		l.UserSearchRoot,
		templateValues{ Username: ldap.EscapeFilter(username) },
		l.UserFilter,
		[]string{ l.EmailField, l.FirstNameField, l.LastNameField },
		conn,
		ctx,
	)
	if err != nil {
		return nil, LDAPError
	}

	if len(result.Entries) == 0 {
		return nil, LDAPNotFoundError
	}

	userEntry := result.Entries[0]
	if err := conn.Bind(userEntry.DN, password); err != nil {
		return nil, AuthenticationError
	}

	return userEntry, nil
}

func (l *ldapUserProvider) getUserLDAPGroupLinks(username, userDN string, conn ldap.Client, ctx context.Context) (ent.GroupLinks, error) {
	response, err := l.ldapSearch(
		l.GroupSearchRoot,
		templateValues{
			Username: ldap.EscapeFilter(username),
			UserDN: userDN,
		},
		l.GroupFilter,
		[]string{l.GroupAttribute},
		conn,
		ctx,
	)
	if err != nil {
		return nil, err
	}

	if len(response.Entries) == 0 {
		return nil, LDAPNotFoundError
	}

	var resourceSpecMatchers []predicate.GroupLink
	groupNames := make([]string, len(response.Entries))
	for i, group := range response.Entries {
		groupNames[i] = group.GetAttributeValue(l.GroupAttribute)
		resourceSpecMatchers = append(resourceSpecMatchers, grouplink.ResourceSpecEQ(groupNames[i]))
	}

	zerolog.Ctx(ctx).Debug().
		Strs("groupsFound", groupNames).
		Msg("Found LDAP groups")

	return ent.FromContext(ctx).GroupLink.Query().
		Where(
			grouplink.And(
				grouplink.TypeEQ("LDAP:" + l.Name),
				grouplink.Or(resourceSpecMatchers...),
			),
		).
		WithClaimGroup(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		All(ctx)
}

func (l *ldapUserProvider) createLocalUser(username string, userEntry *ldap.Entry, ctx context.Context) (*ent.User, error) {
	logger := zerolog.Ctx(ctx)

	fname := userEntry.GetAttributeValue(l.FirstNameField)
	lname := userEntry.GetAttributeValue(l.LastNameField)
	email := userEntry.GetAttributeValue(l.EmailField)

	logger.Debug().
		Str("fname", fname).
		Str("fnameField", l.FirstNameField).
		Str("lname", lname).
		Str("lnameField", l.LastNameField).
		Str("email", email).
		Str("emailField", l.EmailField).
		Str("username", username).
		Msg("Creating User")

	if fname == "" || lname == "" || email == "" {
			return nil, LDAPError
	}

	return ent.FromContext(ctx).User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSource(LDAP_SOURCE).
		Save(ctx)
}

func (l *ldapUserProvider) ldapSearch(searchRoot string, fillTemplate templateValues, filterTemplate *template.Template, attributes []string, conn ldap.Client, ctx context.Context) (*ldap.SearchResult, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("url", l.ServerURL).
		Str("user", fillTemplate.Username).
		Str("userDN", fillTemplate.UserDN).
		Str("template", filterTemplate.Root.String()).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "ldapUserProvider.ldapSearch")
	defer span.End()

	filterBuilder := &strings.Builder{}
	if err := filterTemplate.Execute(filterBuilder, fillTemplate); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not fill filter template")
		return nil, err
	}

	filterString := filterBuilder.String()
	searchRequest := ldap.NewSearchRequest(
		searchRoot,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		l.SearchTimeout,
		false,
		filterString,
		attributes,
		nil,
	)

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("filter", filterString).
		Msg("Executing LDAP search")

	return conn.Search(searchRequest)
}
