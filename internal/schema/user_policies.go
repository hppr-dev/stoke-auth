package schema

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ent"

	"stoke/internal/ent/privacy"

	"github.com/rs/zerolog"
	"hppr.dev/stoke"
)

type superTokenOrSelf struct {}

// Checks whether user who is modifying the token is either the same person or a superuser
// This assumes that the token in the ctx has been validated and verified
func (s superTokenOrSelf) EvalMutation(ctx context.Context, m ent.Mutation) error {
	conf := cfg.Ctx(ctx)
	token := stoke.Token(ctx)
	userM := m.(*ent.UserMutation)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "schema.userMutationPolicy").
		Str("token", token.Raw).
		Stringer("op", userM.Op()).
		Strs("fields", userM.Fields()).
		Logger()

	sc, ok := token.Claims.(*stoke.Claims)
	if !ok {
		logger.Error().Msg("Could not convert claims to stoke claims")
		return privacy.Deny
	}

	mc := sc.StokeClaims
	if super, ok := mc["stk"]; ok && super == "S" {
		logger.Info().Msg("Allowing superuser to modify user.")
		return privacy.Allow
	}

	modUser, ok := userM.Username()
	if !ok {
		logger.Error().Msg("Could not determine user")
		return privacy.Denyf("Could not determine user")
	}

	logger = logger.With().
		Str("user_to_modify", modUser).
		Logger()

	if userClaim, ok := conf.Tokens.UserInfo["username"]; ok {
		if reqUser, ok := mc[userClaim]; ok {
			if reqUser == modUser {
				logger.Info().Msg("User changed.")
				return privacy.Allow
			} else {
				logger.Warn().
					Str("requesting_user", reqUser).
					Msg("User tried to modify")
			}
		}
	}

	return privacy.Denyf("User lacks permission to update %s.", modUser)
}
