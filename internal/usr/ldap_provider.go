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

	"github.com/go-ldap/ldap"
	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

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
	SkipCertificateVerify bool
}

type templateValues struct {
	Username string
	Email    string
	UserDN   string
}

// Init implements Provider.
func (l LDAPUserProvider) Init(context.Context) error {
	return nil
}

// AddUser implements Provider.
// LDAP users should be added upon login
func (l LDAPUserProvider) AddUser(provider ProviderType, fname string, lname string, email string, username string, password string, superUser bool, ctx context.Context) error {
	return ProviderTypeNotSupported
}

// GetUserClaims implements Provider.
func (l LDAPUserProvider) GetUserClaims(username string, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.GetUserClaims")
	defer span.End()

	client := ent.FromContext(ctx)
	tx, err := client.Tx(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not start transaction")
	}

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", username).
		Msg("Getting user claims")

	conn, err := ldap.DialURL(l.ServerURL)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Msg("Could not connect to LDAP server")
		return nil, nil, err
	}
	defer conn.Close()

	if err:= conn.Bind(l.BindUserDN, l.BindUserPassword); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("bindUserDN", l.BindUserDN).
			Msg("Bind user authentication failed")
		return nil, nil, err
	}

	usr, userDN, err := l.getOrCreateUser(username, password, conn, ctx, tx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Msg("Could not get or create user")
		return nil, nil, err
	}

	response, err := l.ldapSearch(
		l.GroupSearchRoot,
		templateValues{ Username: usr.Username, Email: usr.Email, UserDN: userDN },
		l.GroupFilter,
		[]string{l.GroupAttribute},
		conn,
		ctx,
	)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			AnErr("rollbackErr", tx.Rollback()).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Msg("Group search failed")
		return nil, nil, err
	}

	var resourceSpecMatchers []predicate.GroupLink
	groupNames := make([]string, len(response.Entries))
	for i, group := range response.Entries {
		groupNames[i] = group.GetAttributeValue(l.GroupAttribute)
		resourceSpecMatchers = append(resourceSpecMatchers, grouplink.ResourceSpecEQ(groupNames[i]))
	}

	if len(response.Entries) == 0 {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Strs("groupsFound", groupNames).
			AnErr("rollbackErr", tx.Rollback()).
			Msg("Did not find any groups")
		return nil, nil, errors.New("Did not find LDAP groups")
	}

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Strs("groupsFound", groupNames).
		Msg("Found LDAP groups")

	groupLinks, err := tx.GroupLink.Query().
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
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			AnErr("rollbackErr", tx.Rollback()).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Msg("Could not get groupLinks")
		return nil, nil, err
	}

	var claimGroups ent.ClaimGroups
	for _, grouplink := range groupLinks {
		claimGroups = append(claimGroups, grouplink.Edges.ClaimGroup)
	}
	// TODO figure out a way to update local groups with ldap in a better way
	// This should update the local user without producing duplicates?
	if _, err := usr.Update().AddClaimGroups(claimGroups...).Save(ctx); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			AnErr("rollbackErr", tx.Rollback()).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Msg("Failed to add LDAP groups to local user")
		return nil, nil, err
	}

	var claims ent.Claims
	for _, link := range groupLinks {
		claims = append(claims, link.Edges.ClaimGroup.Edges.Claims...)
	}

	return usr, claims, tx.Commit()
}

// UpdateUserPassword implements Provider.
func (l LDAPUserProvider) UpdateUserPassword(provider ProviderType, username string, oldPassword string, newPassword string, force bool, ctx context.Context) error {
	return ProviderTypeNotSupported
}

// Returns a list of claims if there are updates to the user
func (l LDAPUserProvider) getOrCreateUser(username, password string, conn ldap.Client, ctx context.Context, tx *ent.Tx) (*ent.User, string, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.getUser")
	defer span.End()

	result, err := l.ldapSearch(
		l.UserSearchRoot,
		templateValues{ Username: username },
		l.UserFilter,
		[]string{ l.EmailField, l.FirstNameField, l.LastNameField },
		conn,
		ctx,
	)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Msg("User search failed")
		return nil, "", err
	}

	if len(result.Entries) != 1 {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Str("url", l.ServerURL).
			Str("user", username).
			Int("numEntries", len(result.Entries)).
			Msg("Got unexpected number of results")
		return nil, "", errors.New("Got unexpected number of results")
	}

	userEntry := result.Entries[0]
	if err := conn.Bind(userEntry.DN, password); err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("url", l.ServerURL).
			Str("user", username).
			Msg("User authentication failed")
		return nil, "", errors.New("User authentication failed")
	}

	usr, err := tx.User.Query().
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
	if err == nil {
		return usr, userEntry.DN, nil
	} else if !ent.IsNotFound(err) {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", username).
			Str("userDN", userEntry.DN).
			Msg("Error retrieving user")
		return nil, "", err
	}

	logger.Info().
		Func(otelzerolog.AddTracingContext(span)).
		Err(err).
		Str("username", username).
		Msg("Could not pull user locally, creating...")

	fname := userEntry.GetAttributeValue(l.FirstNameField)
	lname := userEntry.GetAttributeValue(l.LastNameField)
	email := userEntry.GetAttributeValue(l.EmailField)

	if fname == "" || lname == "" || email == "" {
		logger.Error().
			Str("fname", fname).
			Str("fnameField", l.FirstNameField).
			Str("lname", lname).
			Str("lnameField", l.LastNameField).
			Str("email", email).
			Str("emailField", l.EmailField).
			Str("username", username).
			Msg("Could not retreive all required attributes")
			return nil, "", errors.New("LDAP config error")
	}

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Err(err).
		Str("fname", fname).
		Str("lname", lname).
		Str("email", email).
		Str("username", username).
		Msg("Creating User.")

	usr, err = tx.User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSource("LDAP").
		Save(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("fname", fname).
			Str("lname", lname).
			Str("email", email).
			Str("username", username).
			Msg("Could not create user")
		return nil, "" , err
	}

	return usr, userEntry.DN, nil
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
			Str("email", fillTemplate.Email).
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
		Str("filterString", filterString).
		Str("template", filterTemplate.Root.String()).
		Msg("Executing LDAP search")

	return conn.Search(searchRequest)
}
