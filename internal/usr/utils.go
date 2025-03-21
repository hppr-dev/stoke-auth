package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"slices"
	"stoke/internal/ent"
	"stoke/internal/ent/schema/policy"
	"stoke/internal/ent/user"
	"golang.org/x/crypto/argon2"
)


func HashPass(pass, salt string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func GenSalt() string {
	saltBytes := make([]byte, 32)
	_, _ = rand.Read(saltBytes)
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

func allUserClaims(u *ent.User) ent.Claims {
	var allClaims ent.Claims
	for _, group := range u.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	return allClaims
}
