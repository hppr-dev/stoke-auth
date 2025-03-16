package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"stoke/internal/ent"
	"stoke/internal/schema/policy"
	"stoke/internal/usr"
)

type Users struct {
	// Enable checking/creating stoke admin per-operation claims
	CreateStokeClaims bool              `json:"create_stoke_claims"`
	// Policy Configuration
	PolicyConfig PolicyConfig           `json:"policy_config"`
	// Configs for providers
	Providers         []*ProviderConfig `json:"providers"`
}

type PolicyConfig struct {
	// Allow superuser override protective policies
	AllowSuperuserOverride bool  `json:"allow_superuser_override"`
	// Whether to disallow any changes to the database after start up
	ReadOnlyMode bool            `json:"read_only_mode"`
	// Users that are not allowed to be changed
	ProtectedUsers []string      `json:"protected_users"`
	// Groups that are not allowed to be changed
	ProtectedGroups []string     `json:"protected_groups"`
	// Claims that are not allowed to be changed
	ProtectedClaims []string     `json:"protected_claims"`
}

func (u Users) withContext(ctx context.Context) context.Context {
	ctx = u.PolicyConfig.withContext(ctx)

	providerList := usr.NewProviderList()

	if u.CreateStokeClaims {
		providerList.CheckCreateForStokeClaims(ctx)
	}

	for _, prov := range u.Providers {
		providerList.AddForeignProvider(prov.Name, prov.CreateProvider(ctx))
	}

	return providerList.WithContext(ctx)
}

func (p PolicyConfig) withContext(ctx context.Context) context.Context {
	conf := Ctx(ctx)
	return policy.ConfigurePolicies(
		p.ProtectedUsers,
		p.ProtectedClaims,
		p.ProtectedGroups,
		conf.Tokens.UserInfo["username"],
		p.ReadOnlyMode,
		p.AllowSuperuserOverride,
		ctx,
	)

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
