package policy

import (
	"context"
	"stoke/internal/ent/privacy"

	"hppr.dev/stoke"
)

type ctxKey byte
const (
	policyConfigCtxKey ctxKey = iota
	allowDBInitKey
)

type policyConfig struct {
	protectedUsernames       []string
	protectedClaimShortNames []string
	protectedGroupNames      []string
	usernameClaim            string
	readOnlyMode             bool
	allowSuperuserOverride   bool
}

func ConfigurePolicies(usernames, claims, groups []string, usernameClaim string, readOnly, allowSuperuserOverride bool, ctx context.Context) context.Context {
	return context.WithValue(ctx, policyConfigCtxKey, &policyConfig{
		protectedUsernames:       usernames,
		protectedClaimShortNames: claims,
		protectedGroupNames:      groups,
		usernameClaim:            usernameClaim,
		readOnlyMode:             readOnly,
		allowSuperuserOverride:   allowSuperuserOverride,
	})
}

func policyFromCtx(ctx context.Context) *policyConfig {
	return ctx.Value(policyConfigCtxKey).(*policyConfig)
}

func isInReadOnlyMode(ctx context.Context) bool {
	return policyFromCtx(ctx).readOnlyMode
}

func getClaimsOrDeny(ctx context.Context) (map[string]string, error) {
	token := stoke.Token(ctx)

	sc, ok := token.Claims.(*stoke.Claims)
	if !ok {
		return nil, privacy.Denyf("Token invalid")
	}

	return sc.StokeClaims, nil
}

func allowChangesBySuperuser(ctx context.Context, claims map[string]string) error {
	if super, ok := claims["stk"]; policyFromCtx(ctx).allowSuperuserOverride && ok && super == "S" {
		return privacy.Allow
	}
	return nil
}


func BypassDatabasePolicies(ctx context.Context) context.Context {
	return context.WithValue(ctx, allowDBInitKey, true)
}

func hasBypassSet(ctx context.Context) bool {
	return ctx.Value(allowDBInitKey) != nil
}
