package usr

import "context"

type providerCtxKey struct {}

func ProviderFromCtx(ctx context.Context) *ProviderList {
	return ctx.Value(providerCtxKey{}).(*ProviderList)
}

func (l *ProviderList) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, providerCtxKey{}, l)
}
