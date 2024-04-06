package usr

import "context"

func ProviderFromCtx(ctx context.Context) Provider {
	return ctx.Value("user-provider").(Provider)
}
