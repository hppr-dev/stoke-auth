package web

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ent/ogent"
	"strings"
)

func (h *entityHandler) AvailableProviders(ctx context.Context) (*ogent.AvailableProvidersOK, error) {
	config := cfg.Ctx(ctx)
	providers := []ogent.AvailableProvidersOKItem{}
	for _, p := range config.Users.Providers {
		providers = append(providers, ogent.AvailableProvidersOKItem{
			Name:         p.Name,
			ProviderType: strings.ToUpper(p.ProviderType),
			TypeSpec:     p.TypeSpec(),
		})
	}
	basePath := strings.TrimRight(config.Server.BasePath, "/")
	return &ogent.AvailableProvidersOK{
		Providers:     providers,
		BaseAdminPath: basePath,
	}, nil
}
