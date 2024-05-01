package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/grouplink"
	"stoke/internal/ent/predicate"
	"stoke/internal/ent/user"
	"stoke/internal/tel"
	"strings"
	"text/template"

	"github.com/go-ldap/ldap/v3"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

var (
	AuthenticationError = errors.New("Could not authenticate user")
	NoLinkedGroupsError = errors.New("No linked groups associated with user")
	LDAPNotFoundError = errors.New("No results")
	LDAPError = errors.New("Error communicating with LDAP")
)

type LDAPConnector interface {
	Connect(string, ...ldap.DialOpt) (ldap.Client, error)
}

type LDAPUserProvider struct {
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

	LocalProvider
}

type templateValues struct {
	Username string
	UserDN   string
}

type ldapDialer struct {}

// Connect using the ldap.DialURL function
func (ldapDialer) Connect(url string, dialOpts ...ldap.DialOpt) (ldap.Client, error) {
	return ldap.DialURL(url, dialOpts...)
}

// Set the ldap connector to use.
// Should only be needed for testing, but could also be used to alter connection behaviour
func (l *LDAPUserProvider) SetConnector(c LDAPConnector) {
	l.connector = c
}

// Set the ldap connector to use the default ldap connector
func (l *LDAPUserProvider) DefaultConnector() {
	l.connector = ldapDialer{}
}

// GetUserClaims looks up claims that are associated in LDAP
func (l LDAPUserProvider) GetUserClaims(username, password string, _ bool, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("url", l.ServerURL).
		Str("username", username).
		Str("bindUserDN", l.BindUserDN).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.GetUserClaims")
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
		return l.LocalProvider.GetUserClaims(username, password, true, ctx)
	}
	defer conn.Close()

	if err:= conn.Bind(l.BindUserDN, l.BindUserPassword); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Bind user authentication failed")
		return l.LocalProvider.GetUserClaims(username, password, true, ctx)
	}

	usr, groupLinks, err := l.getOrCreateUser(username, password, conn, ctx)
	if errors.Is(LDAPNotFoundError, err) || errors.Is(LDAPError, err) {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not find User in ldap")
		return l.LocalProvider.GetUserClaims(username, password, true, ctx)

	} else if errors.Is(NoLinkedGroupsError, err) {
		if usr == nil {
			return nil, nil, AuthenticationError
		}
	} else if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not get or create user")
		return nil, nil, err
	}

	// LDAP user groups that the user already has
	userGroups := usr.Edges.ClaimGroups

	var addClaimGroups ent.ClaimGroups
	var found bool
	for _, grouplink := range groupLinks {
		linkGroup := grouplink.Edges.ClaimGroup
		found = false
		for _, userGroup := range userGroups {
			if userGroup.ID == linkGroup.ID {
				found = true
				break
			}
		}
		if !found {
			addClaimGroups = append(addClaimGroups, linkGroup)
		}
	}

	var delClaimGroups ent.ClaimGroups
	for _, userGroup := range userGroups {
		found = false
		for _, groupLink := range groupLinks {
			if groupLink.Edges.ClaimGroup.ID == userGroup.ID {
				found = true
				break
			}
		}
		if !found {
			delClaimGroups = append(delClaimGroups, userGroup)
		}
	}

	if len(addClaimGroups) > 0 || len(delClaimGroups) > 0 {
		_, err := usr.Update().
			AddClaimGroups(addClaimGroups...).
			RemoveClaimGroups(delClaimGroups...).
			Save(ctx)
		if err != nil {
			logger.Error().
				Func(otelzerolog.AddTracingContext(span)).
				Err(err).
				Msg("Failed to add LDAP groups to local user")
			return nil, nil, err
		}
	}

	return l.LocalProvider.GetUserClaims(username, "", false, ctx)
}

// Creates the user if it exists in LDAP
func (l LDAPUserProvider) getOrCreateUser(username, password string, conn ldap.Client, ctx context.Context) (*ent.User, ent.GroupLinks, error) {
	logger := zerolog.Ctx(ctx).With().
			Str("url", l.ServerURL).
			Str("user", username).
			Logger()

	userEntry, err := l.getLDAPUser(username, password, conn, ctx)
	if err != nil{
		return nil, nil, err
	}

	usr, dbErr := ent.FromContext(ctx).User.Query().
		Where(
			user.And(
				user.Or(
					user.UsernameEQ(username),
					user.EmailEQ(username),
				),
				user.SourceEQ("LDAP"),
			),
		).
		WithClaimGroups(func (q *ent.ClaimGroupQuery) {
			q.Where(
				claimgroup.HasGroupLinksWith(
					grouplink.TypeEQ("LDAP"),
				),
			)
		}).
		Only(ctx)
	if dbErr != nil && !ent.IsNotFound(dbErr) {
		return nil, nil, dbErr
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
		logger.Error().
			Err(err).
			Str("userDN", userEntry.DN).
			Msg("User has no linked groups")
		return usr, groupLinks, NoLinkedGroupsError
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

	logger.Debug().
		Int("numGroups", len(groupLinks)).
		Msg("Found user with group links")

	return usr, groupLinks, nil
}

func (l LDAPUserProvider) getLDAPUser(username, password string, conn ldap.Client, ctx context.Context) (*ldap.Entry, error) {
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

func (l LDAPUserProvider) getUserLDAPGroupLinks(username, userDN string, conn ldap.Client, ctx context.Context) (ent.GroupLinks, error) {
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
				grouplink.TypeEQ("LDAP"),
				grouplink.Or(resourceSpecMatchers...),
			),
		).
		WithClaimGroup(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		All(ctx)
}

func (l LDAPUserProvider) createLocalUser(username string, userEntry *ldap.Entry, ctx context.Context) (*ent.User, error) {
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

	usr, err := ent.FromContext(ctx).User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSource("LDAP").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	return usr, nil
}

func (l LDAPUserProvider) ldapSearch(searchRoot string, fillTemplate templateValues, filterTemplate *template.Template, attributes []string, conn ldap.Client, ctx context.Context) (*ldap.SearchResult, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("url", l.ServerURL).
		Str("user", fillTemplate.Username).
		Str("userDN", fillTemplate.UserDN).
		Str("template", filterTemplate.Root.String()).
		Logger()

	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.ldapSearch")
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
