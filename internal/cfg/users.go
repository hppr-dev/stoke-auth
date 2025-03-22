package cfg

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"stoke/internal/ent"
	"stoke/internal/ent/schema/policy"
	"stoke/internal/usr"
	"strings"

	"github.com/ghodss/yaml"

	"github.com/rs/zerolog"
)

type Users struct {
	// Enable checking/creating stoke admin per-operation claims
	CreateStokeClaims bool              `json:"create_stoke_claims"`
	// Policy Configuration
	PolicyConfig PolicyConfig           `json:"policy_config"`
	// Directory to pull provider definitions from
	ProviderConfigDir string            `json:"provider_config_dir"`
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
	logger := zerolog.Ctx(ctx)
	ctx = u.PolicyConfig.withContext(ctx)

	if u.ProviderConfigDir == "" {
		u.ProviderConfigDir = "/etc/stoke/providers.d/"
	}

	if stat, err := os.Stat(u.ProviderConfigDir); err == nil && stat.IsDir() {
		files, err := os.ReadDir(u.ProviderConfigDir)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Could not read provider config directory")
		}
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".yaml") || strings.HasSuffix(f.Name(), ".yml") {
				provFile, err := os.ReadFile(path.Join(u.ProviderConfigDir, f.Name()))
				if err != nil {
					logger.Error().
						Err(err).
						Str("filename", f.Name()).
						Msg("Could not read provider config file")
						continue
				}
				newProv := &ProviderConfig{}
				if err := yaml.Unmarshal(provFile, newProv); err != nil {
					logger.Error().
						Err(err).
						Str("filename", f.Name()).
						Msg("Could not marshal provider config file")
						continue
				}
				logger.Info().
					Str("filename", f.Name()).
					Msg("Provider config file read")
				u.Providers = append(u.Providers, newProv)
			}
		}
	}

	providerList := usr.NewProviderList()

	if u.CreateStokeClaims {
		if err := providerList.CheckCreateForStokeClaims(ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("Error while creating stoke claims")
		}
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
