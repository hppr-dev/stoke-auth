package policy

import (
	"context"
	"slices"
	"stoke/internal/cfg"
	"stoke/internal/ent"

	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/privacy"

	"github.com/rs/zerolog"
)

type ClaimGroupMutationPolicy struct {}

// Singular policy for ClaimGroup entity:
//   * Allows changes if bypass is set
//   * Disallows changes if read only mode is set
//   * Allow changes by super user
//   * Disallows changes to groups that are in the protected users specified in config
func (p ClaimGroupMutationPolicy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	groupM := m.(*ent.ClaimGroupMutation)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "policy.GroupMutationPolicy").
		Strs("fields", m.Fields()).
		Logger()

	if hasBypassSet(ctx) {
		logger.Info().Msg("Policies Bypassed")
		return privacy.Allow
	}

	if isInReadOnlyMode(ctx) {
		logger.Info().Msg("Server is in read-only mode")
		return privacy.Denyf("Server is running in read-only mode")
	}

	modClaimGroups, err := p.getTargetClaimGroupsOrDeny(ctx, groupM)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get target claim groups")
		return err
	}

	claims, err := getClaimsOrDeny(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get user claims")
		return err
	}

	if err := allowChangesBySuperuser(ctx, claims); err != nil {
		logger.Info().Msg("Superuser allowed update") 
		return err
	}

	if err := p.denyChangesToProtectedEntities(ctx, modClaimGroups); err != nil {
		logger.Warn().Err(err).Msg("Protected entity")
		return err
	}

	logger.Warn().Msg("Allow by default")
	return privacy.Allow
}

func (p ClaimGroupMutationPolicy) getTargetClaimGroupsOrDeny(ctx context.Context, m *ent.ClaimGroupMutation) (ent.ClaimGroups, error) {
	ids, err := m.IDs(ctx)
	if err != nil {
		return nil, privacy.Denyf("Could not determine target group IDs")
	}

	modClaimGroups, err := m.Client().ClaimGroup.Query().Where(claimgroup.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, privacy.Denyf("Could not determine target group IDs")
	}

	return modClaimGroups, nil
}

func (p ClaimGroupMutationPolicy) denyChangesToProtectedEntities(ctx context.Context, groups ent.ClaimGroups) error {
	for _, group := range groups {
		if slices.Contains(cfg.Ctx(ctx).Users.PolicyConfig.ProtectedGroups, group.Name) {
			return privacy.Denyf("ClaimGroup %s is read-only", group.Name)
		}
	}
	return nil
}
