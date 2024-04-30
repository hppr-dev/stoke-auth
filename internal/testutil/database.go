package testutil

import (
	"context"
	"encoding/base64"
	"errors"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/enttest"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)


type DatabaseMutation func(*ent.Client)

func WithDatabase(t *testing.T, mutations ...DatabaseMutation) ContextOption {
	return func(ctx context.Context) context.Context {
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		for _, mut := range mutations {
			mut(client)
		}
		return ent.NewContext(ctx, client)
	}
}

func ForeverKey() DatabaseMutation {
	return KeyWithExpires(time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC))
}

func KeyWithExpires(exp time.Time) DatabaseMutation {
	return func(client *ent.Client) {
		client.PrivateKey.Create().
			SetExpires(exp).
			SetText("DHGQKw0oDDcMcZArDSgMNwxxkCsNKAw3DHGQKw0oDDe1+1s+xW4vzlPSPGN3OTEStdBKaW3SHjMRGJL5rk6IAA==").
			SaveX(context.Background())
	}
}

func ForeverKeyWithText(text string) DatabaseMutation {
	return func(client *ent.Client) {
		client.PrivateKey.Create().
			SetExpires(time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)).
			SetText(text).
			SaveX(context.Background())
	}
}

func DefaultUser() DatabaseMutation {
	return func(client *ent.Client) {
		User(
			UserInfo("frank", "lash", "flash", "flash@hppr.dev"),
			Source("LOCAL"),
			Password("flashpass"),
			Group(
				GroupInfo("Speeders", "People with super speed"),
				Claim(
					ClaimInfo("power", "pow", "speed", "Has super speed"),
				),
			),
		)(client)
	}
}

type UserOption func(*ent.UserCreate)

func User(opts ...UserOption) DatabaseMutation {
	return func(client *ent.Client) {
		userCreate := client.User.Create()
		for _, opt := range opts {
			opt(userCreate)
		}
		userCreate.SaveX(context.Background())
	}
}

func UserInfo(fname, lname, username, email string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetFname(fname)
		u.SetLname(lname)
		u.SetUsername(username)
		u.SetEmail(email)
	}
}

func Source(src string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetSource(src)
	}
}

func Password(pass string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetSalt("HELLOWORLD")
		u.SetPassword(HashPass(pass, "HELLOWORLD"))
	}
}

func GroupFromName(name string) UserOption {
	return func(u *ent.UserCreate) {
		group := u.Mutation().Client().ClaimGroup.Query().Where(claimgroup.NameEQ(name)).FirstX(context.Background())
		u.AddClaimGroups(group)
	}
}

type GroupOption func(*ent.ClaimGroupCreate)

func Group(opts ...GroupOption) UserOption {
	return func(u *ent.UserCreate) {
		groupCreate := u.Mutation().Client().ClaimGroup.Create()
		for _, opt := range opts {
			opt(groupCreate)
		}
		group := groupCreate.SaveX(context.Background())
		u.AddClaimGroups(group)
	}
}

func GroupInfo(name, desc string) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		c.SetName(name)
		c.SetDescription(desc)
	}
}

func ClaimFromName(name string) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		claim := c.Mutation().Client().Claim.Query().Where(claim.NameEQ(name)).FirstX(context.Background())
		c.AddClaims(claim)
	}
}

type ClaimOption func(*ent.ClaimCreate)

func Claim(opts ...ClaimOption) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		claimCreate := c.Mutation().Client().Claim.Create()
		for _, opt := range opts {
			opt(claimCreate)
		}
		claim := claimCreate.SaveX(context.Background())
		c.AddClaims(claim)
	}
}

func ClaimInfo(name, shortName, value, desc string) ClaimOption {
	return func (c *ent.ClaimCreate) {
		c.SetDescription(desc)
		c.SetShortName(shortName)
		c.SetName(name)
		c.SetValue(value)
	}
}

func HashPass(pass, salt string) string {
		return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

type entErrHelper struct {}
var simEntErr = errors.New("Simulated ent error") 
func (entErrHelper) Query(context.Context, ent.Query) (ent.Value, error) { return nil, simEntErr }
func (entErrHelper) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) { return nil, simEntErr }

func ReturnsMutateErrors() DatabaseMutation {
	return func(client *ent.Client) {
		client.Use(func(ent.Mutator) ent.Mutator{
			return entErrHelper{}
		})
	}
}

func ReturnsReadErrors() DatabaseMutation {
	return func(client *ent.Client) {
		client.Intercept(ent.InterceptFunc(func(ent.Querier) ent.Querier {
			return entErrHelper{}
		}))
	}
}

func ReturnsAllErrors() DatabaseMutation {
	return func(client *ent.Client) {
		ReturnsReadErrors()(client)
		ReturnsMutateErrors()(client)
	}
}
