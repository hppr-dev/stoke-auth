package ctx

import (
	"context"
	"errors"
	"stoke/internal/cfg"
	"stoke/internal/ent"
	"stoke/internal/key"
	"stoke/internal/usr"

	"github.com/rs/zerolog"
)

type ShutdownFunc func(Context) error
type InitFunc func() error

type Context struct {
	Config cfg.Config
	Issuer key.TokenIssuer
	UserProvider usr.Provider
	DB *ent.Client
	RootLogger zerolog.Logger
	AppContext context.Context
	
	startupFuncs  []InitFunc
	shutdownFuncs []ShutdownFunc
}

func (c *Context) OnStartup(f InitFunc) {
	c.startupFuncs = append(c.startupFuncs, f)
}

func (c *Context) OnShutdown(f ShutdownFunc) {
	c.shutdownFuncs = append(c.shutdownFuncs, f)
}

func (c Context) Startup() error {
	var err error
	for _, f := range c.startupFuncs {
		err = errors.Join(err, f())
	}
	return err
}

func (c Context) GracefulShutdown() error {
	var err error
	for _, f := range c.shutdownFuncs {
		err = errors.Join(err, f(c))
	}
	return err
}
