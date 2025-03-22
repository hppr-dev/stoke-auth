package web

import (
	"context"
	"stoke/internal/cfg"
	"stoke/internal/ent/ogent"
	"strings"
)

func (h *entityHandler) AvailableProviders(ctx context.Context) ([]ogent.AvailableProvidersOKItem, error) {
	config := cfg.Ctx(ctx)
	res := []ogent.AvailableProvidersOKItem{}
	for _, p := range config.Users.Providers {
		res = append(res, ogent.AvailableProvidersOKItem{
			Name:         p.Name,
			ProviderType: strings.ToUpper(p.ProviderType),
			TypeSpec:     p.TypeSpec(),
		})
	}
	return res, nil
}
