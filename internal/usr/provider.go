package usr

import (
	"context"
	"errors"
	"stoke/internal/ent"

	"github.com/rs/zerolog"
)

type ProviderType int

const (
	LOCAL ProviderType = iota
	LDAP
)

var ProviderTypeNotSupported = errors.New("Provider type not supported")

type Provider interface {
	Init(context.Context) error
	GetUserClaims(user, password string, ctx context.Context) (*ent.User, ent.Claims, error)
  AddUser(provider ProviderType, fname, lname, email, username, password string, superUser bool, ctx context.Context) error
	// Only need to handle this on local providers. Not sure if we need/want this on other providers
  UpdateUserPassword(provider ProviderType, username, oldPassword, newPassword string, force bool, ctx context.Context) error
}

type MultiProvider struct {
	providers map[ProviderType]Provider
}

func (m *MultiProvider) Add(t ProviderType, p Provider) {
	if m.providers == nil {
		m.providers = make(map[ProviderType]Provider)
	}
	m.providers[t] = p
}

func (m *MultiProvider) Init(ctx context.Context) error {
	zerolog.Ctx(ctx).Info().
		Int("numProviders", len(m.providers)).
		Msg("Initializing multiprovider...")

	for _, p := range m.providers {
		err := p.Init(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *MultiProvider) AddUser(provider ProviderType, fname, lname, email, username, password string, superUser bool, ctx context.Context) error {
	p, ok := m.providers[provider]
	if !ok {
		return ProviderTypeNotSupported
	}
	return p.AddUser(provider, fname, lname, email, username, password, superUser, ctx)
}

func (m *MultiProvider) UpdateUserPassword(provider ProviderType, username, oldPassword, newPassword string, force bool, ctx context.Context) error {
	if provider != LOCAL {
		return ProviderTypeNotSupported
	}
	return m.providers[LOCAL].UpdateUserPassword(provider, username, oldPassword, newPassword, force, ctx)
}

func (m *MultiProvider) GetUserClaims(username, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	p, ok := m.providers[LDAP]
	if ok {
		return p.GetUserClaims(username, password, ctx)
	}
	return m.providers[LOCAL].GetUserClaims(username, password, ctx)
}
