package key

import (
	"context"
)

func IssuerFromCtx(ctx context.Context) TokenIssuer {
	return ctx.Value("issuer").(TokenIssuer)
}
