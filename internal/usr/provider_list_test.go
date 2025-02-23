package usr_test

import (
	"context"
	"errors"
	"stoke/internal/ent"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/user"
	tu "stoke/internal/testutil"
	"stoke/internal/usr"
	"testing"
)


func TestNoForeignProviders(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "hello", "hello@local"),
				tu.Password("world"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("Admin Claim", "adm", "Y", "Administrator Claim"),
					),
				),
			),
		),
	)

	emptyProviderList := usr.NewProviderList()

	user, claims, err := emptyProviderList.GetUserClaims("hello", "world", ctx)
	if err != nil {
		t.Fatalf("Failed to get user claims from local: %v", err)
	}

	if user == nil {
		t.Fatal("Returned user was nil")
	}

	if len(claims) != 1 {
		t.Fatalf("Returned claims were not as expected: %v", claims)
	}

	c := claims[0]

	if c.ShortName != "adm" || c.Name != "Admin Claim" || c.Value != "Y" {
		t.Fatalf("Returned claim was not as expected: %v", c)
	}
}

type MockProvider struct {
	AddGroup    string
	RemoveGroup string
	ReturnValue error
}

func (m *MockProvider) UpdateUserClaims(username, _ string, ctx context.Context) error {
	foundUser, _ := ent.FromContext(ctx).User.Query().
		Where(user.UsernameEQ(username)).
		First(ctx)
	if foundUser != nil {
		if m.AddGroup != "" {
			foundGroup, _ := ent.FromContext(ctx).ClaimGroup.Query().
				Where(claimgroup.NameEQ(m.AddGroup)).
				First(ctx)
			foundUser.Update().AddClaimGroups(foundGroup).SaveX(ctx)
		}
		if m.RemoveGroup != "" {
			foundGroup, _ := ent.FromContext(ctx).ClaimGroup.Query().
				Where(claimgroup.NameEQ(m.RemoveGroup)).
				First(ctx)
			foundUser.Update().RemoveClaimGroups(foundGroup).SaveX(ctx)
		}
	}
	return m.ReturnValue
}

func TestReturnsForeignInfoWhenForeignDoesNotReturnAnError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("foreign", "user", "user1", "hello@local"),
				tu.Source("CUSTOM"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("Admin Claim", "adm", "Y", "Administrator Claim"),
					),
				),
			),
		),
	)

	pl := usr.NewProviderList()

	foreign := &MockProvider{}
	pl.AddForeignProvider(foreign)

	user, claims, err := pl.GetUserClaims("user1", "chooch", ctx)
	if err != nil {
		t.Fatalf("GetUserClaims returned an error: %v", err)
	}

	if user == nil || user.Username != "user1" {
		t.Fatalf("User was not as expected: %v", user)
	}

	if len(claims) != 1 {
		t.Fatalf("Claims length was not as expected: %v", claims)
	}

	if claims[0].ShortName != "adm" || claims[0].Value != "Y" {
		t.Fatalf("Returned claim was not as expected: %v", claims[0])
	}
}

func TestReturnsLocalInfoWhenForeignReturnsAuthSourceError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "local@local"),
				tu.Password("plop"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("Admin Claim", "adm", "Y", "Administrator Claim"),
					),
				),
			),
			tu.User(
				tu.UserInfo("foreign", "user", "user1", "hello@local"),
				tu.Source("CUSTOM"),
				tu.GroupFromName("admin group"),
			),
		),
	)

	pl := usr.NewProviderList()

	foreign := &MockProvider{
		ReturnValue: usr.AuthSourceError,
	}
	pl.AddForeignProvider(foreign)

	user, claims, err := pl.GetUserClaims("localuser", "plop", ctx)
	if err != nil {
		t.Fatalf("GetUserClaims returned an error: %v", err)
	}

	if user == nil || user.Username != "localuser" {
		t.Fatalf("User was not as expected: %v", user)
	}

	if len(claims) != 1 {
		t.Fatalf("Claims length was not as expected: %v", claims)
	}

	if claims[0].ShortName != "adm" || claims[0].Value != "Y" {
		t.Fatalf("Returned claim was not as expected: %v", claims[0])
	}
}

func TestReturnsForeignInfoWithMultipleForeignProviders(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "local@local"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("success group", "this group should be added when the user is successfully found"),
					tu.Claim(
						tu.ClaimInfo("Hello", "hel", "wor", "Hello world claim"),
					),
				),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("Admin Claim", "adm", "Y", "Administrator Claim"),
					),
				),
			),
			tu.User(
				tu.UserInfo("foreign", "user", "user1", "hello@local"),
				tu.Source("CUSTOM"),
				// No Groups to start
			),
		),
	)

	pl := usr.NewProviderList()

	foreignFail := &MockProvider{
		ReturnValue: usr.AuthSourceError,
	}
	foreignSuccess := &MockProvider{
		AddGroup : "success group",
	}

	pl.AddForeignProvider(foreignFail)
	pl.AddForeignProvider(foreignSuccess)


	user, claims, err := pl.GetUserClaims("user1", "somepass", ctx)
	if err != nil {
		t.Fatalf("GetUserClaims returned an error: %v", err)
	}

	if user == nil || user.Username != "user1" {
		t.Fatalf("User was not as expected: %v", user)
	}

	if len(claims) != 1 {
		t.Fatalf("Claims length was not as expected: %v", claims)
	}

	if claims[0].ShortName != "hel" || claims[0].Value != "wor" {
		t.Fatalf("Returned claim was not as expected: %v", claims[0])
	}
}

func TestReturnsBadPasswordWhenAProviderReturnsAuthenticationError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "local@local"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("Admin Claim", "adm", "Y", "Administrator Claim"),
					),
				),
			),
		),
	)

	pl := usr.NewProviderList()

	foreignFail := &MockProvider{
		ReturnValue: usr.AuthenticationError,
	}
	pl.AddForeignProvider(foreignFail)

	_, _, err := pl.GetUserClaims("user1", "somepass", ctx)
	if err == nil {
		t.Fatal("GetUserClaims did not return an error")
	}

	if !errors.Is(err, usr.AuthenticationError) {
		t.Fatal("Did not return authentication error")
	}
}
