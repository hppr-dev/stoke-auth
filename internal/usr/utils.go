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

func findGroupChanges(u *ent.User, groupLinks ent.GroupLinks) (ent.ClaimGroups, ent.ClaimGroups) {
	userGroups := u.Edges.ClaimGroups

	var addClaimGroups ent.ClaimGroups
	var found bool
	for _, grouplink := range groupLinks {
		linkGroup := grouplink.Edges.ClaimGroup
		found = false
		for _, userGroup := range userGroups {
			if userGroup.ID == linkGroup.ID {
				found = true
				break
			}
		}
		if !found {
			addClaimGroups = append(addClaimGroups, linkGroup)
		}
	}

	var delClaimGroups ent.ClaimGroups
	for _, userGroup := range userGroups {
		found = false
		for _, groupLink := range groupLinks {
			if groupLink.Edges.ClaimGroup.ID == userGroup.ID {
				found = true
				break
			}
		}
		if !found {
			delClaimGroups = append(delClaimGroups, userGroup)
		}
	}

	return addClaimGroups, delClaimGroups
}

func applyGroupChanges(add, del ent.ClaimGroups, u *ent.User, ctx context.Context) error {
	var err error
	if len(add) > 0 || len(del) > 0 {
		_, err = u.Update().
			AddClaimGroups(add...).
			RemoveClaimGroups(del...).
			Save(ctx)
	}
	return err
}

func retreiveLocalUser(username string, ctx context.Context) (*ent.User, error) {
	return ent.FromContext(ctx).User.Query().
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

	u, err := retreiveLocalUser(username, ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Msg("Could not find user")
		return nil, nil, err
	}

	allClaims := allUserClaims(u)
	logger.Debug().
		Str("username", username).
		Func(otelzerolog.AddTracingContext(span)).
		Interface("claims", allClaims).
		Msg("Claims found")
	return u, allClaims, nil
}

func allUserClaims(u *ent.User) ent.Claims {
	var allClaims ent.Claims
	for _, group := range u.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	return allClaims
}
