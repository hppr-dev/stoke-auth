package policy

import (
	"context"
	"slices"
	"stoke/internal/cfg"
	"stoke/internal/ent"

	"stoke/internal/ent/privacy"

	"github.com/rs/zerolog"
)

type UserMutationPolicy struct {}

// Singular policy for User entity:
//   * Disallows changes if read only mode is set
//   * Disallows changes to users that are in the protected users specified in config
//   * Disallows changes to multiple users at once
//   * Allows the logged in user to modify themselves
//   * Allows superusers to modify anyone else
func (p UserMutationPolicy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	userM := m.(*ent.UserMutation)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "policy.UserMutationPolicy").
		Strs("fields", m.Fields()).
		Logger()

	if hasBypassSet(ctx) {
		logger.Info().Msg("Policies Bypassed")
		return privacy.Allow
	}

	if isInReadOnlyMode(ctx) {
		logger.Warn().Msg("Server is in read-only mode")
		return privacy.Denyf("Server is running in read-only mode")
	}

	user, err := p.getTargetUserOrDeny(ctx, userM)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get target user")
		return err
	}

	if err := p.denyChangesToProtectedEntities(ctx, user); err != nil {
		logger.Warn().Err(err).Msg("Protected entity")
		return err
	}

	claims, err := getClaimsOrDeny(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get user claims")
		return err
	}

	if err := p.allowChangesToSelf(ctx, user, claims); err != nil {
		logger.Info().Err(err).Msg("Self allowed update")
		return err
	}

	if err := allowChangesBySuperuser(ctx, claims); err != nil {
		logger.Info().Msg("Superuser allowed update") 
		return err
	}

	logger.Error().Msg("Allow by default")
	return privacy.Allow
}

func (p UserMutationPolicy) getTargetUserOrDeny(ctx context.Context, m *ent.UserMutation) (*ent.User, error) {
	ids, err := m.IDs(ctx)
	if err != nil {
		return nil, privacy.Denyf("Could not determine target user")
	}
	if len(ids) != 1 {
		return nil, privacy.Denyf("Can not change more than one user at once")
	}

	modUser, err := m.Client().User.Get(ctx, ids[0])
	if err != nil {
		return nil, privacy.Denyf("Could not determine user")
	}

	return modUser, nil
}

func (p UserMutationPolicy) denyChangesToProtectedEntities(ctx context.Context, user *ent.User) error {
	if slices.Contains(cfg.Ctx(ctx).Users.PolicyConfig.ProtectedUsers, user.Username) {
		return privacy.Denyf("User %s is read-only", user.Username)
	}
	return nil
}

func (p UserMutationPolicy) allowChangesToSelf(ctx context.Context, user *ent.User, claims map[string]string) error {
	usernameClaim, ok := cfg.Ctx(ctx).Tokens.UserInfo["username"]
	if !ok {
		return privacy.Denyf("Could not determine username")
	}
	username, _ := claims[usernameClaim]
	if user.Username == username {
		return privacy.Allow
	}
	return nil
}
