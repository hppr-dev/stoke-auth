package policy

import (
	"context"
	"slices"
	"stoke/internal/ent"

	"stoke/internal/ent/claimgroup"
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
//   * Disallows assigning superuser priviledges from non-superusers
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

	if userM.Op().Is(ent.OpCreate) {
		logger.Debug().Msg("Allowing user creation")
		return privacy.Allow
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

	if err := p.denyAssigningSuperuserFromNonSuperuser(ctx, userM, claims); err != nil {
		logger.Warn().Err(err).Msg("Non-superuser trying to assign super priviledges")
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
	if slices.Contains(policyFromCtx(ctx).protectedUsernames, user.Username) {
		return privacy.Denyf("User %s is read-only", user.Username)
	}
	return nil
}

func (p UserMutationPolicy) allowChangesToSelf(ctx context.Context, user *ent.User, claims map[string]string) error {
	username := claims[policyFromCtx(ctx).usernameClaim]
	if user.Username == username {
		return privacy.Allow
	}
	return nil
}

func (p UserMutationPolicy) denyAssigningSuperuserFromNonSuperuser(ctx context.Context, userM *ent.UserMutation, claims map[string]string) error {
	if slices.Contains(userM.AddedEdges(), "claim_groups") {
		groupIDs := userM.ClaimGroupsIDs()
		mutGroups, err := userM.Client().ClaimGroup.Query().
			Where(
				claimgroup.IDIn(groupIDs...),
			).
			WithClaims().
			All(ctx)
		if err != nil {
			return privacy.Denyf("Could not determine assigned groups")
		}

		userIsNotSuper := claims["stk"] != "S"
		for _, group := range mutGroups {
			for _, claim := range group.Edges.Claims {
				if userIsNotSuper && claim.ShortName == "stk" && claim.Value == "S" {
					return privacy.Denyf("Cannot assign superuser claim from non superuser")
				}
			}
		}
	}
	return nil
}
