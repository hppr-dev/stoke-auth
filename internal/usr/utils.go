package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"stoke/internal/ent"
	"stoke/internal/ent/user"
	"stoke/internal/tel"

	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"golang.org/x/crypto/argon2"
)


func HashPass(pass, salt string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func GenSalt() string {
	saltBytes := make([]byte, 32)
	rand.Read(saltBytes)
	return base64.StdEncoding.EncodeToString(saltBytes)
}

// retreives the user from the local database. If the user exists, it returns the claims that are associated
func retreiveLocalClaims(username string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "usr.retreiveLocalClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", username).
		Msg("Getting user claims")

	usr, err := ent.FromContext(ctx).User.Query().
		Where(
			user.And(
				user.Or(
					user.UsernameEQ(username),
					user.EmailEQ(username),
				),
			),
		).
		WithClaimGroups(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		Only(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Msg("Could not find user")
		return nil, nil, err
	}

	var allClaims ent.Claims
	for _, group := range usr.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	logger.Debug().
		Str("username", username).
		Func(otelzerolog.AddTracingContext(span)).
		Interface("claims", allClaims).
		Msg("Claims found")
	return usr, allClaims, nil
}
