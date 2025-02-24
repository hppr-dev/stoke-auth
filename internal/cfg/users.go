package cfg

import (
	"context"
	"stoke/internal/usr"
)

type foreignProvider interface {
	UpdateUserClaims(username, password string, ctx context.Context) error
}

type ProviderConfig interface {
	CreateProvider(context.Context) foreignProvider
}

type Users struct {
	// Enable checking/creating stoke admin per-operation claims
	CreateStokeClaims bool             `json:"create_stoke_claims"`
	// Configs for providers
	Providers         []ProviderConfig `json:"providers"`

}

func (u Users) withContext(ctx context.Context) context.Context {
	providerList := usr.NewProviderList()

	if u.CreateStokeClaims {
		providerList.CheckCreateForStokeClaims(ctx)
	}

	return providerList.WithContext(ctx)
}
