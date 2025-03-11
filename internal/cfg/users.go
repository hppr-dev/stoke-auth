package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"stoke/internal/ent"
	"stoke/internal/usr"
)

type Users struct {
	// Enable checking/creating stoke admin per-operation claims
	CreateStokeClaims bool             `json:"create_stoke_claims"`
	// Configs for providers
	Providers         []*ProviderConfig `json:"providers"`
}

func (u Users) withContext(ctx context.Context) context.Context {
	providerList := usr.NewProviderList()

	if u.CreateStokeClaims {
		providerList.CheckCreateForStokeClaims(ctx)
	}

	for _, prov := range u.Providers {
		providerList.AddForeignProvider(prov.Name, prov.CreateProvider(ctx))
	}

	return providerList.WithContext(ctx)
}

type ProviderConfig struct {
	providerConfig
	ProviderType   string `json:"type"`
	Name           string `json:"name"`
}

type providerConfig interface {
	CreateProvider(context.Context) foreignProvider
	TypeSpec() string
}

type foreignProvider interface {
	UpdateUserClaims(username, password string, ctx context.Context) (*ent.User, error)
}

func (pc *ProviderConfig) UnmarshalJSON(b []byte) error {
	temp := struct {
		ProviderType string `json:"type"`
		Name         string `json:"name"`
	}{}
	err := json.Unmarshal(b, &temp)
	if err != nil {
		return err
	}

	pc.ProviderType = temp.ProviderType
	pc.Name = temp.Name
	switch(pc.ProviderType) {
	case "ldap", "LDAP":
		pc.providerConfig = &LDAPProviderConfig{}
		return json.Unmarshal(b, pc.providerConfig)
	case "oidc", "OIDC":
		pc.providerConfig = &OIDCProviderConfig{}
		return json.Unmarshal(b, pc.providerConfig)
	}
	return fmt.Errorf("Provider type not supported: %s", temp.ProviderType)
}
