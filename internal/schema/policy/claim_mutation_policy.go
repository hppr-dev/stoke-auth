package policy

import (
	"context"
	"slices"
	"stoke/internal/cfg"
	"stoke/internal/ent"

	"stoke/internal/ent/claim"
	"stoke/internal/ent/privacy"

	"github.com/rs/zerolog"
)

type ClaimMutationPolicy struct {}

// Singular policy for Claim entity:
//   * Allows changes if bypass is set
//   * Disallows changes if read only mode is set
//   * Allow changes by super user
//   * Disallows changes to claims that are in the protected users specified in config
func (p ClaimMutationPolicy) EvalMutation(ctx context.Context, m ent.Mutation) error {
	claimM := m.(*ent.ClaimMutation)

	logger := zerolog.Ctx(ctx).With().
		Str("component", "policy.ClaimMutationPolicy").
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

	modClaims, err := p.getTargetClaimsOrDeny(ctx, claimM)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get target claims")
		return err
	}

	claims, err := getClaimsOrDeny(ctx)
	if err != nil {
		logger.Warn().Err(err).Msg("Could not get user claims")
		return err
	}

	if err := allowChangesBySuperuser(ctx, claims); err != nil {
		logger.Info().Err(err).Msg("Superuser allowed update.")
		return err
	}

	if err := p.denyChangesToProtectedEntities(ctx, modClaims); err != nil {
		logger.Warn().Err(err).Msg("Protected entity")
		return err
	}

	logger.Warn().Msg("Allow by default")
	return privacy.Allow
}

func (p ClaimMutationPolicy) getTargetClaimsOrDeny(ctx context.Context, m *ent.ClaimMutation) (ent.Claims, error) {
	ids, err := m.IDs(ctx)
	if err != nil {
		return nil, privacy.Denyf("Could not determine target claim IDs")
	}

	modClaims, err := m.Client().Claim.Query().Where(claim.IDIn(ids...)).All(ctx)
	if err != nil {
		return nil, privacy.Denyf("Could not determine target claim IDs")
	}

	return modClaims, nil
}

func (p ClaimMutationPolicy) denyChangesToProtectedEntities(ctx context.Context, claims ent.Claims) error {
	for _, claim := range claims {
		if slices.Contains(cfg.Ctx(ctx).Users.PolicyConfig.ProtectedClaims, claim.ShortName) {
			return privacy.Denyf("Claim %s is read-only", claim.ShortName)
		}
	}
	return nil
}
