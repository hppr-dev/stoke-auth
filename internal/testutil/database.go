package testutil

import (
	"context"
	"encoding/base64"
	"errors"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/enttest"
	"stoke/internal/ent/schema/policy"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/argon2"
)

var bypassCtx = policy.BypassDatabasePolicies(context.Background())

type DatabaseMutation func(*ent.Client)

// Adds a mock database to the context
func WithDatabase(t *testing.T, mutations ...DatabaseMutation) ContextOption {
	return func(ctx context.Context) context.Context {
		ctx = policy.BypassDatabasePolicies(ctx)
		client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&_fk=1")
		for _, mut := range mutations {
			mut(client)
		}
		return ent.NewContext(ctx, client)
	}
}

// Creates a key that doesn't expire until year 5000
func ForeverKey() DatabaseMutation {
	return KeyWithExpires(time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC))
}

// Creates a key that expires at a given time
func KeyWithExpires(exp time.Time) DatabaseMutation {
	return func(client *ent.Client) {
		client.PrivateKey.Create().
			SetExpires(exp).
			SetText("DHGQKw0oDDcMcZArDSgMNwxxkCsNKAw3DHGQKw0oDDe1-1s-xW4vzlPSPGN3OTEStdBKaW3SHjMRGJL5rk6IAA==").
			SaveX(bypassCtx)
	}
}

// Creates a key that doesn't expire until year 5000 with the given key text
func ForeverKeyWithText(text string) DatabaseMutation {
	return func(client *ent.Client) {
		client.PrivateKey.Create().
			SetExpires(time.Date(5000, time.January, 10, 10, 10, 10, 10, time.UTC)).
			SetText(text).
			SaveX(bypassCtx)
	}
}

// Creates a default user in the database to use with tests
// fname: frank, lname: lash, username: flash, email: flash@hppr.dev
// source: LOCAL, password: flashpass
// Groups:
//   * Speeders -> People with super speed 
//     * power (pow=speed) -> Has super speed 
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

// Creates a user with the given options.
// This is the entrypoint to insert users/groups/claims into the mock database.
// All other user related entities are created in relation to user entities
func User(opts ...UserOption) DatabaseMutation {
	return func(client *ent.Client) {
		userCreate := client.User.Create()
		for _, opt := range opts {
			opt(userCreate)
		}
		userCreate.SaveX(bypassCtx)
	}
}

// Sets the users information
func UserInfo(fname, lname, username, email string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetFname(fname)
		u.SetLname(lname)
		u.SetUsername(username)
		u.SetEmail(email)
	}
}

// Sets the user's source field
func Source(src string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetSource(src)
	}
}

// Sets the user password. Hashes the password and sets the salt to HELLOWORLD
func Password(pass string) UserOption {
	return func(u *ent.UserCreate) {
		u.SetSalt("HELLOWORLD")
		u.SetPassword(HashPass(pass, "HELLOWORLD"))
	}
}

// Lookup a group by name and add it to the user. The group should already be created
func GroupFromName(name string) UserOption {
	return func(u *ent.UserCreate) {
		group := u.Mutation().Client().ClaimGroup.Query().Where(claimgroup.NameEQ(name)).FirstX(bypassCtx)
		u.AddClaimGroups(group)
	}
}

type GroupOption func(*ent.ClaimGroupCreate)

// Creates a group with given options and add it to the user.
// Should have GroupInfo and at least 1 Claim
func Group(opts ...GroupOption) UserOption {
	return func(u *ent.UserCreate) {
		groupCreate := u.Mutation().Client().ClaimGroup.Create()
		for _, opt := range opts {
			opt(groupCreate)
		}
		group := groupCreate.SaveX(bypassCtx)
		u.AddClaimGroups(group)
	}
}

// Sets group information
func GroupInfo(name, desc string) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		c.SetName(name)
		c.SetDescription(desc)
	}
}

// Creates an LDAP group link and adds it to the group
func LDAPLink(providerName, groupName string) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		link := c.Mutation().Client().GroupLink.Create().
			SetResourceSpec(groupName).
			SetType("LDAP:" + providerName).
			SaveX(bypassCtx)
		c.AddGroupLinks(link)
	}
}

// Add a claim to the group using a name to look up. The claim should be created before calling this.
func ClaimFromName(name string) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		claim := c.Mutation().Client().Claim.Query().Where(claim.NameEQ(name)).FirstX(bypassCtx)
		c.AddClaims(claim)
	}
}

type ClaimOption func(*ent.ClaimCreate)

// Create a claim with given options, i.e. Claim(ClaimInfo("name", "short", "value", "description")
func Claim(opts ...ClaimOption) GroupOption {
	return func(c *ent.ClaimGroupCreate) {
		claimCreate := c.Mutation().Client().Claim.Create()
		for _, opt := range opts {
			opt(claimCreate)
		}
		claim := claimCreate.SaveX(bypassCtx)
		c.AddClaims(claim)
	}
}

// Set Claim information
func ClaimInfo(name, shortName, value, desc string) ClaimOption {
	return func (c *ent.ClaimCreate) {
		c.SetDescription(desc)
		c.SetShortName(shortName)
		c.SetName(name)
		c.SetValue(value)
	}
}

// Hashes passwords as they are hashed by the local user provider.
// This is duplicated here to not have test depend on server code
func HashPass(pass, salt string) string {
		return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

type entErrHelper struct {}
var simEntErr = errors.New("Simulated ent error") 
func (entErrHelper) Query(context.Context, ent.Query) (ent.Value, error) { return nil, simEntErr }
func (entErrHelper) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) { return nil, simEntErr }

// Makes database mutations return an error. All tables are affected by default
func ReturnsMutateErrors(tables ...string) DatabaseMutation {
	errFunc := func(ent.Mutator) ent.Mutator { return entErrHelper{} }
	return func(client *ent.Client) {
		if len(tables) == 0 {
			client.Use(errFunc)
			return
		}
		for _, t := range tables {
			switch t {
			case "claim":
				client.Claim.Use(errFunc)
			case "claimgroup", "claim-group", "claimGroup":
				client.ClaimGroup.Use(errFunc)
			case "user":
				client.User.Use(errFunc)
			case "grouplink", "group-link", "groupLink":
				client.GroupLink.Use(errFunc)
			case "privatekey", "private-key", "privateKey":
				client.PrivateKey.Use(errFunc)
			default:
				panic("Unknown table for mutate errors: " + t)
			}
		}
	}
}

// Makes database reads return an error. All tables are affected by default
func ReturnsReadErrors(tables ...string) DatabaseMutation {
	errFunc := ent.InterceptFunc(func(ent.Querier) ent.Querier { return entErrHelper{} })
	return func(client *ent.Client) {
		if len(tables) == 0 {
			client.Intercept(errFunc)
			return
		}
		for _, t := range tables {
			switch t {
			case "claim":
				client.Claim.Intercept(errFunc)
			case "claimgroup", "claim-group", "claimGroup":
				client.ClaimGroup.Intercept(errFunc)
			case "user":
				client.User.Intercept(errFunc)
			case "grouplink", "group-link", "groupLink":
				client.GroupLink.Intercept(errFunc)
			case "privatekey", "private-key", "privateKey":
				client.PrivateKey.Intercept(errFunc)
			default:
				panic("Unknown table for mutate errors: " + t)
			}
		}
	}
}

// Makes all database operations return errors. All tables are affected by default
func ReturnsAllErrors(tables ...string) DatabaseMutation {
	return func(client *ent.Client) {
		ReturnsReadErrors(tables...)(client)
		ReturnsMutateErrors(tables...)(client)
	}
}
