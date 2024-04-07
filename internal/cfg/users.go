package cfg

import (
	"context"
	"stoke/internal/usr"

	"github.com/rs/zerolog"
)

type Users struct {
	// Nothing in here until I get LDAP working
}

func (u Users) withContext(ctx context.Context) context.Context {
	// TODO need to be able to configure which provider to use
	// Will always have user
	localUserProvider := usr.LocalProvider{}
	err := localUserProvider.Init(ctx)
	if err != nil {
		zerolog.Ctx(ctx).
			Fatal().
			Err(err).
			Msg("Could not initialize local user database")
	}

	return context.WithValue(ctx, "user-provider", localUserProvider)
}
