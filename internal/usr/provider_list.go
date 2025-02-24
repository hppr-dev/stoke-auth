package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"
)

type provider interface {
	// Updates user claims in the local database if and only if we are able successfully authenticate the user
	// Returns AuthenticationError if the password is bad
	// Returns AuthSourceError if there is an issue with pulling from the authentication source
	UpdateUserClaims(user, password string, ctx context.Context) error
}

type ProviderList struct {
	*localProvider
	foreignProviders []provider
}

func NewProviderList() *ProviderList {
	return &ProviderList{
		foreignProviders: []provider{},
		localProvider: &localProvider{},
	}
}

// Gets the given users claims
// Updates the local database with claims received from foreignProviders
// If the user does not exist in any foreignProviders, the localProvider is checked
// If the foreignProviders fail to produce claims, local claims are given, if and only if the user is a local user
// Claims are tracked in the local database regardless of which provider the claims were derived from
// Claims may only be pulled from a single provider at a time
func (p *ProviderList) GetUserClaims(user, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	for _, provider := range p.foreignProviders {
		if err := provider.UpdateUserClaims(user, password, ctx); err == nil {
			break
		} else if errors.Is(err, AuthenticationError) {
			return nil, nil, err
		}
	}
	return p.localProvider.GetUserClaims(user, password, ctx)
}

func (p *ProviderList) AddForeignProvider(newProvider provider) {
	p.foreignProviders = append(p.foreignProviders, newProvider)
}
