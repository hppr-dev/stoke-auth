package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"

	"github.com/rs/zerolog"
)

type provider interface {
	// Updates user claims in the local database if and only if we are able successfully authenticate the user
	// Returns AuthenticationError if the password is bad
	// Returns AuthSourceError if there is an issue with pulling from the authentication source
	UpdateUserClaims(user, password string, ctx context.Context) (*ent.User, error)
}

type ProviderList struct {
	*localProvider
	foreignProviders map[string]provider
}

func NewProviderList() *ProviderList {
	return &ProviderList{
		foreignProviders: make(map[string]provider),
		localProvider: &localProvider{},
	}
}

// Gets the given users claims
// Updates the local database with claims received from foreignProviders
// If the user does not exist in any foreignProviders, the localProvider is checked
// If the foreignProviders fail to produce claims, local claims are given, if and only if the user is a local user
// Claims are tracked in the local database regardless of which provider the claims were derived from
// Claims may only be pulled from a single provider at a time
func (p *ProviderList) GetUserClaims(username, password, providerID string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "usr.ProviderList").
		Str("username", username).
		Str("provider", providerID).
		Logger()

	var u *ent.User
	var err error

	prov, found := p.foreignProviders[providerID]
	if providerID == "" && len(p.foreignProviders) == 1 {
		for _, v := range p.foreignProviders {
			prov = v
			found = true
		}
	}

	if found {
		u, err = prov.UpdateUserClaims(username, password, ctx);
		logger.Debug().
			Err(err).
			Interface("user", u).
			Msg("Done checking provider")
		if errors.Is(err, AuthenticationError) {
			return nil, nil, err
		}
	}

	return p.localProvider.GetUserClaims(username, password, u, ctx)
}

func (p *ProviderList) AddForeignProvider(name string, newProvider provider) {
	p.foreignProviders[name] = newProvider
}
