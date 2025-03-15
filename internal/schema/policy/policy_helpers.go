package policy

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ent/privacy"

	"hppr.dev/stoke"
)

func isInReadOnlyMode(ctx context.Context) bool {
	return cfg.Ctx(ctx).Users.PolicyConfig.ReadOnlyMode
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
	if super, ok := claims["stk"]; cfg.Ctx(ctx).Users.PolicyConfig.AllowSuperuserOverride && ok && super == "S" {
		return privacy.Allow
	}
	return nil
}

type allowDBInitKey struct {}

func BypassDatabasePolicies(ctx context.Context) context.Context {
	return context.WithValue(ctx, allowDBInitKey{}, true)
}

func hasBypassSet(ctx context.Context) bool {
	return ctx.Value(allowDBInitKey{}) != nil
}
