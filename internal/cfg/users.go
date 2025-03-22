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
	// Directory to pull provider definitions from
	ProviderConfigDir string            `json:"provider_config_dir"`
	// Directory to pull user database init files from
	UserInitDir       string            `json:"user_init_dir"`
	// Single file to initialize database from
	UserInitFile      string            `json:"user_init_file"`
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

func (u *Users) withContext(ctx context.Context) context.Context {
	logger := zerolog.Ctx(ctx)
	ctx = u.PolicyConfig.withContext(ctx)

	if u.ProviderConfigDir == "" {
		u.ProviderConfigDir = "/etc/stoke/providers.d/"
	}

	if err := u.initLocalDatabase(ctx); err != nil {
		logger.Warn().
			Err(err).
			Msg("Could not init local database")
	}

	if err := u.parseProviders(ctx); err != nil {
		logger.Warn().
			Err(err).
			Msg("Error parsing providers")
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

func (u *Users) parseProviders(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().
		Str("provider_config_dir", u.ProviderConfigDir).
		Logger()

	logger.Debug().Msg("Parsing providers")
	if stat, err := os.Stat(u.ProviderConfigDir); err == nil && stat.IsDir() {
		files, err := os.ReadDir(u.ProviderConfigDir)
		if err != nil {
			logger.Error().
				Err(err).
				Msg("Could not read provider config directory")
			return err
		}

		for _, f := range files {
			if isYAMLFile(f.Name()) {
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
	return nil
}

func (u *Users) initLocalDatabase(ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().
		Str("user_init_file", u.UserInitFile).
		Str("user_init_dir", u.UserInitDir).
		Logger()

	logger.Debug().Msg("Initializing local database")
	if u.UserInitFile != "" {
		if err := InitializeDatabaseFromFile(u.UserInitFile, ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("Could not read init file")
			return err
		}
		logger.Info().Msg("Initialized database from file")
	}
	if u.UserInitDir != "" {
		if stat, err := os.Stat(u.UserInitDir); err == nil && stat.IsDir() {
			files, err := os.ReadDir(u.UserInitDir)
			if err != nil {
				logger.Error().
					Err(err).
					Msg("Could not read init directory")
				return err
			}
			for _, f := range files {
				if isYAMLFile(f.Name()) {
					if err := InitializeDatabaseFromFile(path.Join(u.UserInitDir, f.Name()), ctx); err != nil{
						logger.Error().
							Err(err).
							Str("filename", f.Name()).
							Msg("Could not init from file")
						continue
					}
					logger.Info().
						Str("filename", f.Name()).
						Msg("Initialized database from file")
				}
			}
		}
	}
	return nil
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

func isYAMLFile(filename string) bool {
	return strings.HasSuffix(filename, ".yaml") || strings.HasSuffix(filename, ".yml")
}
