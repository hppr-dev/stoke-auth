package ctx

import (
	"stoke/internal/cfg"
	"stoke/internal/key"
)

type Context struct {
	Config cfg.Config
	Issuer key.TokenIssuer
}
