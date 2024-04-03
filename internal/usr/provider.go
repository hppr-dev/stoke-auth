package usr

import (
	"errors"
	"stoke/internal/ent"
)

type ProviderType int

const (
	LOCAL ProviderType = iota
	LDAP
)

var ProviderTypeNotSupported = errors.New("Provider type not supported")

type Provider interface {
	Init() error
	GetUserClaims(user, pass string) (*ent.User, ent.Claims, error)
  AddUser(provider ProviderType, fname, lname, email, username, password string, superUser bool) error
  UpdateUser(provider ProviderType, fname, lname, email, username, password string) error
}

type MultiProvider struct {
	providers map[ProviderType]Provider
}

func (m MultiProvider) Add(t ProviderType, p Provider) {
	m.providers[t] = p
}

func (m MultiProvider) Init() error {
	logger.Info().
		Int("numProviders", len(m.providers)).
		Msg("Initializing multiprovider...")

	for _, p := range m.providers {
		err := p.Init()
		if err != nil {
			return err
		}
	}
	return nil
}

func (m MultiProvider) AddUser(provider ProviderType, fname, lname, email, username, password string, superUser bool) error {
	p, ok := m.providers[provider]
	if !ok {
		return ProviderTypeNotSupported
	}
	return p.AddUser(provider, fname, lname, email, username, password, superUser)
}

func (m MultiProvider) UpdateUser(provider ProviderType, fname, lname, email, username, password string) error {
	p, ok := m.providers[provider]
	if !ok {
		return ProviderTypeNotSupported
	}
	return p.UpdateUser(provider, fname, lname, email, username, password)
}

func (m MultiProvider) GetUserClaims(username, password string) (*ent.User, ent.Claims, error) {
	var claims ent.Claims
	var user *ent.User
	for _, p := range m.providers {
		provUser, provClaims, _ := p.GetUserClaims(username, password)
		claims = append(claims, provClaims...)
		if provUser != nil {
			user = provUser
		}
	}
	if len(claims) == 0 {
		logger.Debug().
			Str("username", username).
			Msg("No claims found")
		return nil, nil, errors.New("No claims found")
	}
	return user, claims, nil
}
