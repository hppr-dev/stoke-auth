package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"
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

type LDAPUserProvider struct {
	localProvider         LocalProvider
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
}

type templateValues struct {
	Username string
	UserDN   string
}

var AuthenticationError = errors.New("Could not authenticate user")
var LDAPNotFoundError = errors.New("No results")
var LDAPError = errors.New("Error communicating with LDAP")

// Init implements Provider.
func (l LDAPUserProvider) Init(ctx context.Context) error {
	return l.localProvider.Init(ctx)
}

// AddUser implements Provider.
// LDAP users should be added upon login
func (l LDAPUserProvider) AddUser(fname, lname, email, username, password string, superuser bool, ctx context.Context) error {
	return l.localProvider.AddUser(fname, lname, email, username, password, superuser, ctx)
}

// GetUserClaims implements Provider.
func (l LDAPUserProvider) GetUserClaims(username, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.GetUserClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", username).
		Msg("Getting user claims")

	conn, err := ldap.DialURL(l.ServerURL, l.DialOpts...)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Msg("Could not connect to LDAP server")
		return l.localProvider.GetUserClaims(username, password, ctx)
	}
	defer conn.Close()

	if err:= conn.Bind(l.BindUserDN, l.BindUserPassword); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("bindUserDN", l.BindUserDN).
			Msg("Bind user authentication failed")
		return l.localProvider.GetUserClaims(username, password, ctx)
	}

	usr, groupLinks, err := l.getOrCreateUser(username, password, conn, ctx)
	if errors.Is(LDAPNotFoundError, err) {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("url", l.ServerURL).
			Str("username", username).
			Msg("User not found in ldap")
		return l.localProvider.GetUserClaims(username, password, ctx)

	} else if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Msg("Could not get or create user")
		return nil, nil, err
	}

	var claimGroups ent.ClaimGroups
	for _, grouplink := range groupLinks {
		claimGroups = append(claimGroups, grouplink.Edges.ClaimGroup)
	}

	// TODO implement LDAP group removal if user is no longer a member
	if _, err := usr.Update().AddClaimGroups(claimGroups...).Save(ctx); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Msg("Failed to add LDAP groups to local user")
		return nil, nil, err
	}

	return l.localProvider.GetUserClaims(username, "", ctx)
}

// UpdateUserPassword implements Provider.
// Changes local passwords only. LDAP password change is not supported
func (l LDAPUserProvider) UpdateUserPassword(username, oldPassword, newPassword string, force bool, ctx context.Context) error {
	return l.localProvider.UpdateUserPassword(username, oldPassword, newPassword, force, ctx)
}

// Creates the user if it exists in LDAP
func (l LDAPUserProvider) getOrCreateUser(username, password string, conn ldap.Client, ctx context.Context) (*ent.User, ent.GroupLinks, error) {
	logger := zerolog.Ctx(ctx)

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
		Only(ctx)
	if dbErr != nil && !ent.IsNotFound(dbErr) {
		return nil, nil, err
	}

	groupLinks, err := l.getUserLDAPGroupLinks(username, userEntry.DN, conn, ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Str("userDN", userEntry.DN).
			Msg("Could not get users linked LDAP groups")
		return nil, nil, err
	}

	if len(groupLinks) == 0 {
		logger.Error().
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Str("userDN", userEntry.DN).
			Msg("User has no linked groups")
		return nil, nil, AuthenticationError
	}

	if ent.IsNotFound(dbErr) {
		usr, err = l.createLocalUser(username, userEntry, ctx)
		if err != nil {
			logger.Error().
				Err(err).
				Str("username", username).
				Str("userDN", userEntry.DN).
				Msg("Could not create local user")
			return nil, nil, err
		}
	}

	logger.Debug().
		Str("username", username).
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
		return nil, errors.Join(LDAPError, err)
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
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.ldapSearch")
	defer span.End()

	filterBuilder := &strings.Builder{}
	if err := filterTemplate.Execute(filterBuilder, fillTemplate); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", fillTemplate.Username).
			Str("userDN", fillTemplate.UserDN).
			Str("template", filterTemplate.Root.String()).
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
		Str("template", filterTemplate.Root.String()).
		Msg("Executing LDAP search")

	return conn.Search(searchRequest)
}
