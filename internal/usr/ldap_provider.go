package usr

import (
	"context"
	"stoke/internal/ent"
	"stoke/internal/ent/claimgroup"
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
	SkipCertificateVerify bool
	GroupSearchRoot       string
	GroupFilter           *template.Template
	GroupAttribute        string
	AutoAddGroups         bool
}

// Init implements Provider.
func (l LDAPUserProvider) Init(context.Context) error {
	return nil
}

// Creates the user in the local database as an LDAP user and adds LDAP groups that have been linked (if password is given and correct)
// AddUser implements Provider.
func (l LDAPUserProvider) AddUser(provider ProviderType, fname string, lname string, email string, username string, password string, superUser bool, ctx context.Context) error {
	if provider != LDAP {
		return ProviderTypeNotSupported
	}

	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.AddUser")
	defer span.End()

	usr, err := ent.FromContext(ctx).User.Create().
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
		return err
	}

	// TODO check password/ldap groups before creating local user.
	if password != "" {
		if _, err := l.authenticate(usr, password, true, ctx); err != nil {
			logger.Error().
				Func(otelzerolog.AddTracingContext(span)).
				Err(err).
				Str("fname", fname).
				Str("lname", lname).
				Str("email", email).
				Str("username", username).
				Msg("Could not authenticate user")
			return err
		}
	}

	return nil
}

// GetUserClaims implements Provider.
func (l LDAPUserProvider) GetUserClaims(username string, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.GetUserClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", username).
		Msg("Getting user claims")

	usr, err := ent.FromContext(ctx).User.Query().
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
			q.WithClaims()
			// Local groups only
			q.Where(
				claimgroup.Not(
					claimgroup.HasGroupLinks(),
				),
			)
		}).
		Only(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Msg("Could not find user")
		return nil, nil, err
	}

	claims, err := l.authenticate(usr, password, false, ctx) 
	if err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Msg("User authentication failed.")
		return nil, nil, err
	}

	for _, group := range usr.Edges.ClaimGroups {
		claims = append(claims, group.Edges.Claims...)
	}
	return usr, claims, nil
}

// UpdateUserPassword implements Provider.
func (l LDAPUserProvider) UpdateUserPassword(provider ProviderType, username string, oldPassword string, newPassword string, force bool, ctx context.Context) error {
	return ProviderTypeNotSupported
}

// Returns a list of claims if there are updates to the user
func (l LDAPUserProvider) authenticate(usr *ent.User, password string, updateGroups bool, ctx context.Context) (ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LDAPUserProvider.authenticate")
	defer span.End()

	conn, err := ldap.DialURL(l.ServerURL)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Msg("Could not connect to LDAP server")
		return nil, err
	}
	defer conn.Close()

	if err:= conn.Bind(usr.Username, password); err != nil {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Msg("User authentication failed")
		return nil, err
	}

	groupFilterBuilder := &strings.Builder{}
	if err := l.GroupFilter.Execute(groupFilterBuilder, usr); err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Str("template", l.GroupFilter.Root.String()).
			Msg("Could not fill group filter template")
		return nil, err
	}

	groupFilter := groupFilterBuilder.String()
	groupSearch := ldap.NewSearchRequest(
		l.GroupSearchRoot,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		groupFilter,
		[]string{l.GroupAttribute},
		nil,
	)

	response, err := conn.Search(groupSearch)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("url", l.ServerURL).
			Str("user", usr.Username).
			Str("search", groupFilter).
			Msg("Group search failed")
		return nil, err
	}

	var resourceSpecMatchers []predicate.GroupLink
	groupNames := make([]string, len(response.Entries))
	for i, group := range response.Entries {
		groupNames[i] = group.GetAttributeValue(l.GroupAttribute)
		resourceSpecMatchers = append(resourceSpecMatchers, grouplink.ResourceSpecEQ(groupNames[i]))
	}

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Strs("groupsFound", groupNames).
		Msg("Found LDAP groups")

	groupLinks, err := ent.FromContext(ctx).GroupLink.Query().
		Where(
			grouplink.And(
				grouplink.TypeEQ("LDAP"),
				grouplink.Or(resourceSpecMatchers...),
			),
		).
		WithClaimGroups(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		All(ctx)

	var claims ent.Claims
	for _, link := range groupLinks {
		claims = append(claims, link.Edges.ClaimGroups.Edges.Claims...)
	}

	if l.AutoAddGroups || updateGroups {
		var claimGroups ent.ClaimGroups
		for _, grouplink := range groupLinks {
			// TODO  update ClaimGroups to be ClaimGroup or remove unique
			claimGroups = append(claimGroups, grouplink.Edges.ClaimGroups)
		}
		// This should update the local user without producing duplicates?
		if _, err := usr.Update().AddClaimGroups(claimGroups...).Save(ctx); err != nil {
			logger.Error().
				Func(otelzerolog.AddTracingContext(span)).
				Err(err).
				Str("url", l.ServerURL).
				Str("user", usr.Username).
				Msg("Failed to add LDAP groups to local user")
			return nil, err
		}
	}

	return claims, nil
}
