package ctx

import (
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/key"
	"stoke/internal/usr"

	"github.com/rs/zerolog"
)

type Context struct {
	Config cfg.Config
	Issuer key.TokenIssuer
	UserProvider usr.Provider
	DB *ent.Client
	RootLogger zerolog.Logger
}
