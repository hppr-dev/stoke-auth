package schema

import (
	"context"
	"slices"
	"stoke/internal/cfg"
	"stoke/internal/ent"

	"stoke/internal/ent/privacy"

	"github.com/rs/zerolog"
	"hppr.dev/stoke"
)

type RestrictUpdates struct {
	EntityType  string
	FieldName   string
}

// Restricts updates to users, claims, or groups depending on config
func (r RestrictUpdates) EvalMutation(ctx context.Context, m ent.Mutation) error {
	conf := cfg.Ctx(ctx)
	token := stoke.Token(ctx)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "schema.RestrictUpdates").
		Str("entity_name", r.EntityType).
		Logger()

	var protectList []string
	switch r.EntityType {
	case "users":
		protectList = conf.Users.PolicyConfig.ProtectedUsers
	case "claims":
		protectList = conf.Users.PolicyConfig.ProtectedClaims
	case "groups":
		protectList = conf.Users.PolicyConfig.ProtectedGroups
	default:
		return privacy.Denyf("Unknown restriction entity type: %s", r.EntityType)
	}


	sc, ok := token.Claims.(*stoke.Claims)
	if !ok {
		logger.Error().Msg("Could not convert claims to stoke claims")
		return privacy.Deny
	}

	mc := sc.StokeClaims
	if super, ok := mc["stk"]; conf.Users.PolicyConfig.AllowSuperuserOverride && ok && super == "S" {
		logger.Info().Msg("Allowing superuser")
		return privacy.Deny
	}

	if f, ok := m.Field(r.FieldName); ok {
		val := f.(string)
		if slices.Contains(protectList, val) {
			logger.Info().
				Str("protected_entity", val).
				Msg("Protecting entity")
			return privacy.Denyf("Resource is read-only")
		}
	}

	return privacy.Skip
}
