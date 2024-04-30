package usr_test

import (
	"stoke/internal/ent"
	"stoke/internal/ent/user"
	tu "stoke/internal/testutil"
	"stoke/internal/usr"
	"testing"
)

func TestLocalGetUserClaimsHappy(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("Gordan", "Ramsey", "gramsey", "gramsey@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Password("thisdoesn'ttastegood"),
				tu.Group(
					tu.GroupInfo("Show Hosts", "Hosts of kitchen nightmares"),
					tu.Claim(
						tu.ClaimInfo("Set Access", "set", "allow", "Allows people on to set"),
					),
					tu.Claim(
						tu.ClaimInfo("Kitchen Access", "kit", "deny", "Blocks people from the kitchen"),
					),
				),
			),
		),
	)

	localProvider := usr.LocalProvider{}

	user, claims, err := localProvider.GetUserClaims("gramsey", "thisdoesn'ttastegood", ctx)
	if err != nil {
		t.Fatalf("GetUserClaims returned an error: %v", err)
	}

	if len(claims) != 2 {
		t.Logf("Claims did not match: %v", claims)
		t.Fail()
	}

	for _, claim := range claims {
		switch claim.Name {
		case "Set Access":
			if claim.ShortName != "set" || claim.Value != "allow" {
				t.Logf("Claim did not match: %v", claim)
				t.Fail()
			}
		case "Kitchen Access":
			if claim.ShortName != "kit" || claim.Value != "deny" {
				t.Logf("Claim did not match: %v", claim)
				t.Fail()
			}
		default:
			t.Logf("Received unknown claim: %v", claim)
			t.Fail()
		}
	}
	
	if user.Email != "gramsey@hppr.dev" || user.Fname != "Gordan" ||
			user.Lname != "Ramsey" || user.Username != "gramsey" || user.Source != "LOCAL" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

func TestLocalAddUserHappy(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))

	localProvider := usr.LocalProvider{}

	localProvider.AddUser("Lucas", "Sky", "lsky@hppr.dev", "lsky", "fortsbe", false, ctx)

	client := ent.FromContext(ctx)
	allUsers := client.User.Query().AllX(ctx)

	if len(allUsers) != 1 {
		t.Fatalf("Number of users did not match: %v", allUsers)
	}

	user := allUsers[0]

	if user.Username != "lsky" || user.Source != "LOCAL" ||
			user.Email != "lsky@hppr.dev" || user.Fname != "Lucas" ||
			user.Lname != "Sky" || user.Password == "fortsbe" || user.Salt == "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

func TestLocalUpdateUserPasswordHappy(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("Gordan", "Ramsey", "gramsey", "gramsey@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Password("changeme"),
				tu.Group(
					tu.GroupInfo("Show Hosts", "Hosts of kitchen nightmares"),
					tu.Claim(
						tu.ClaimInfo("Set Access", "set", "allow", "Allows people on to set"),
					),
					tu.Claim(
						tu.ClaimInfo("Kitchen Access", "kit", "deny", "Blocks people from the kitchen"),
					),
				),
			),
		),
	)
	localProvider := usr.LocalProvider{}

	if err := localProvider.UpdateUserPassword("gramsey", "changeme", "somethingelse", false, ctx); err != nil {
		t.Fatalf("Failed to UpdateUserPassword: %v", err)
	}

	user := ent.FromContext(ctx).User.Query().Where(user.UsernameEQ("gramsey")).FirstX(ctx)

	if tu.HashPass("somethingelse", user.Salt) != user.Password {
		t.Fatalf("Changing password did not result in a matching hashed password: %v", user)
	}

	if user.Username != "gramsey" || user.Source != "LOCAL" ||
			user.Email != "gramsey@hppr.dev" || user.Fname != "Gordan" ||
			user.Lname != "Ramsey" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

func TestLocalCheckCreateForSuperUserHappy(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))

	localProvider := usr.LocalProvider{}

	if err := localProvider.CheckCreateForSuperUser(ctx); err != nil {
		t.Fatalf("Failed to CheckCreateForSuperUser: %v", err)
	}

	client := ent.FromContext(ctx)
	allUsers := client.User.Query().AllX(ctx)
	allClaims := client.Claim.Query().AllX(ctx)
	allGroups := client.ClaimGroup.Query().AllX(ctx)

	if len(allUsers) != 1 {
		t.Fatalf("Number of users did not match: %v", allUsers)
	}

	if len(allClaims) != 1 {
		t.Fatalf("Number of claims did not match: %v", allClaims)
	}

	if len(allGroups) != 1 {
		t.Fatalf("Number of groups did not match: %v", allGroups)
	}

	user := allUsers[0]
	claim := allClaims[0]
	group := allGroups[0]

	if user.Username != "sadmin" || user.Source != "LOCAL" ||
			user.Password == "" || user.Salt == "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}

	if claim.ShortName != "srol" || claim.Value != "spr" {
		t.Logf("Claim short name and value did not match: %v", claim)
		t.Fail()
	}

	if group.Name != "Stoke Superusers" {
		t.Logf("Group name did not match: %v", group)
		t.Fail()
	}
}
