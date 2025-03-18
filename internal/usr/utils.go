package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"slices"
	"stoke/internal/ent"
	"stoke/internal/ent/schema/policy"
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

func findGroupChanges(u *ent.User, groupLinks ent.GroupLinks, groupType string) (ent.ClaimGroups, ent.ClaimGroups) {
	userGroups := slices.DeleteFunc(
		slices.Clone(u.Edges.ClaimGroups),
		func (g *ent.ClaimGroup) bool {
			return !slices.ContainsFunc(
				g.Edges.GroupLinks,
				func(l *ent.GroupLink) bool {
					return l.Type == groupType
				},
			)
		},
	)
	linkedGroups := make(ent.ClaimGroups, len(groupLinks))

	for i, link := range groupLinks {
		linkedGroups[i] = link.Edges.ClaimGroup
	}


	// Need to add all groups that are in linkedGroups and not in userGroups
	addClaimGroups := slices.DeleteFunc(
		slices.Clone(linkedGroups),
		func (g1 *ent.ClaimGroup) bool {
			return slices.ContainsFunc(
				userGroups,
				func (g2 *ent.ClaimGroup) bool {
					return g1.ID == g2.ID
				},
			)
		},
	)
	// Need to remove all groups that are in userGroups, but not in linked groups
	delClaimGroups := slices.DeleteFunc(
		slices.Clone(userGroups),
		func (g1 *ent.ClaimGroup) bool {
			return slices.ContainsFunc(
				linkedGroups,
				func (g2 *ent.ClaimGroup) bool {
					return g1.ID == g2.ID
				},
			)
		},
	)

	return addClaimGroups, delClaimGroups
}

func applyGroupChanges(add, del ent.ClaimGroups, u *ent.User, ctx context.Context) (*ent.User, error) {
	ctx = policy.BypassDatabasePolicies(ctx)
	builder := u.Update()
	if len(add) > 0 {
		builder.AddClaimGroups(add...)
	}
	if len(del) > 0 {
		builder.RemoveClaimGroups(del...)
	}
	return builder.Save(ctx)
}

func retreiveLocalUser(username string, ctx context.Context) (*ent.User, error) {
	return ent.FromContext(ctx).User.Query().
		Where(
			user.Or(
				user.UsernameEQ(username),
				user.EmailEQ(username),
			),
		).
		WithClaimGroups(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
			q.WithGroupLinks()
		}).
		Only(ctx)
}

// retreives the user from the local database. If the user exists, it returns the claims that are associated
func retreiveLocalClaims(username string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx).With().
		Str("component", "retreiveLocalClaims").
		Str("username", username).
		Logger()
	ctx, span := tel.GetTracer().Start(ctx, "usr.retreiveLocalClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Msg("Getting user claims")

	u, err := retreiveLocalUser(username, ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Msg("Could not find user")
		return nil, nil, err
	}

	allClaims := allUserClaims(u)
	logger.Debug().
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
