package usr

import (
	"stoke/internal/ent"
	"stoke/internal/ent/user"
	tu "stoke/internal/testutil"
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

	localProvider := localProvider{}

	user, claims, err := localProvider.GetUserClaims("gramsey", "thisdoesn'ttastegood", nil, ctx)
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

func TestGetUserClaimsQueryFailure(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsReadErrors()))

	localProvider := localProvider{}

	if _, _, err := localProvider.GetUserClaims("gramsey", "thisdoesn'ttastegood", nil, ctx); err == nil {
		t.Fatal("GetUserClaims did not return an error")
	}
}

func TestGetUserClaimsBadPassword(t *testing.T) {
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

	localProvider := localProvider{}

	if _, _, err := localProvider.GetUserClaims("gramsey", "yummyintummy", nil, ctx); err == nil {
		t.Fatal("GetUserClaims did not return an error")
	}
}

func TestLocalAddUserHappy(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))

	localProvider := localProvider{}

	if err := localProvider.AddUser("Lucas", "Sky", "lsky@hppr.dev", "lsky", "fortsbe", ctx); err != nil {
		t.Fatalf("AddUser returned an error: %v", err)
	}

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

func TestAddUserReturnsErrorWhenDatabaseFails(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsMutateErrors()))

	localProvider := localProvider{}

	if err := localProvider.AddUser("Lucas", "Sky", "lsky@hppr.dev", "lsky", "fortsbe", ctx); err == nil {
		t.Fatalf("AddUser did not return an error: %v", err)
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
	localProvider := localProvider{}

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

func TestLocalUpdateUserPasswordBadPassword(t *testing.T) {
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
	localProvider := localProvider{}

	if err := localProvider.UpdateUserPassword("gramsey", "dontchangeme", "somethingelse", false, ctx); err == nil {
		t.Log("Did not return error")
		t.Fail()
	}

	user := ent.FromContext(ctx).User.Query().Where(user.UsernameEQ("gramsey")).FirstX(ctx)

	if tu.HashPass("changeme", user.Salt) != user.Password {
		t.Logf("Password changed: %v", user)
		t.Fail()
	}

	if user.Username != "gramsey" || user.Source != "LOCAL" ||
			user.Email != "gramsey@hppr.dev" || user.Fname != "Gordan" ||
			user.Lname != "Ramsey" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

func TestLocalUpdateUserPasswordDatabaseFailure(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsReadErrors()))
	localProvider := localProvider{}

	if err := localProvider.UpdateUserPassword("gramsey", "dontchangeme", "somethingelse", false, ctx); err == nil {
		t.Log("Did not return error")
		t.Fail()
	}
}

func TestLocalCheckCreateForSuperUserHappy(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))

	localProvider := localProvider{}

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

	if claim.ShortName != "stk" || claim.Value != "S" {
		t.Logf("Claim short name and value did not match: %v", claim)
		t.Fail()
	}

	if group.Name != "Stoke Superusers" {
		t.Logf("Group name did not match: %v", group)
		t.Fail()
	}
}

func TestLocalCheckCreateForSuperUserDatabaseClaimFailure(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsMutateErrors("claim")))

	localProvider := localProvider{}

	if err := localProvider.CheckCreateForSuperUser(ctx); err == nil {
		t.Fatalf("Did not return error: %v", err)
	}
}

func TestLocalCheckCreateForSuperUserDatabaseGroupWriteFailure(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsMutateErrors("claimgroup")))

	localProvider := localProvider{}

	if err := localProvider.CheckCreateForSuperUser(ctx); err == nil {
		t.Fatalf("Did not return error: %v", err)
	}
}

func TestLocalCheckCreateForSuperUserDatabaseGroupReadFailure(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t, tu.ReturnsReadErrors("claimgroup")))

	localProvider := localProvider{}

	if err := localProvider.CheckCreateForSuperUser(ctx); err == nil {
		t.Fatalf("Did not return error: %v", err)
	}
}
