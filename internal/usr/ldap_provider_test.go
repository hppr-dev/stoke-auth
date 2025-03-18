package usr_test

import (
	"context"
	"errors"
	"fmt"
	"stoke/internal/ent"
	"stoke/internal/ent/user"
	tu "stoke/internal/testutil"
	"stoke/internal/usr"
	"testing"
	"text/template"
)

type LDAPUserProvider interface {
	SetConnector(usr.LDAPConnector)
	UpdateUserClaims(username, password string, ctx context.Context) (*ent.User, error)
}

// User does not exist in the local database yet
// User has claims from linked groups
func TestLDAPUpdateUserClaimsHappy(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "user@local"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
		),
		tu.WithToken(t),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err != nil {
		t.Fatalf("Failed to get user claims: %v", err)
	}

	user, claims := getUserAndClaims("ldapuser", ctx)

	if len(claims) != 2 {
		t.Logf("Claims did not match: %v", claims)
		t.Fail()
	}

	for _, c := range claims {
		switch c.ShortName {
		case "adm":
			if c.Value != "yes" || c.Name != "admin" {
				t.Logf("Claim did not match: %v", c)
				t.Fail()
			}
		case "usr":
			if c.Value != "yes" || c.Name != "user" {
				t.Logf("Claim did not match: %v", c)
				t.Fail()
			}
		default:
			t.Logf("Got unexpected claim: %v", c)
		}
	}

	if user.Username != "ldapuser" || user.Fname != "ldap" ||
			user.Lname != "user" || user.Email != "luser@hppr.dev" ||
			user.Source != "LDAP" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// User is not in the database and has no linked groups
func TestLDAPUpdateUserClaimsNoLinkedGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "user@local"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if u, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatalf("Was able to get claims for user with no linked groups: %v", u)
	}
}

// User is in the database but they have a linked group that is no longer linked to their LDAP groups
func TestLDAPUpdateUserClaimsRemovedGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("main_ldap", "other_ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "other_ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	_, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx)
	if err != nil {
		t.Fatalf("Failed to get user claims: %v", err)
	}

	user, claims := getUserAndClaims("ldapuser", ctx)

	if len(claims) != 1 {
		t.Logf("Claims length did not match: %v", claims)
		t.Fail()
	}

	c := claims[0]

	if c.ShortName != "usr" || c.Value != "yes" || c.Name != "user" {
		t.Logf("Claim did not match: %v", c)
		t.Fail()
	}

	if user.Username != "ldapuser" || user.Fname != "ldap" ||
			user.Lname != "user" || user.Email != "luser@hppr.dev" ||
			user.Source != "LDAP:main_ldap" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// User is in the database but they have a linked group that is no longer linked to their LDAP groups
// But the database cannot be updated
func TestLDAPUpdateUserClaimsRemovedGroupsReturnsErrorDatabaseFailure(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("main_ldap", "other_ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
			tu.ReturnsMutateErrors("user"),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "other_ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("Did not return an error")
	}
}

// User is in the database but the database cannot be read
func TestLDAPUpdateUserClaimsReturnsErrorDatabaseFailure(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
			),
			tu.ReturnsReadErrors("user"),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("Did not return an error")
	}
}

// User is in the database but all linked groups were removed
func TestLDAPUpdateUserClaimsRemovedAllGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "deleted_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("main_ldap", "deleted_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	_, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx)
	if err == nil {
		t.Fatalf("UpdateUserClaims did not return an error")
	}

	user, claims := getUserAndClaims("ldapuser", ctx)

	if len(claims) != 0 {
		t.Logf("User still has claims: %v", claims)
		t.Fail()
	}

	if len(user.Edges.ClaimGroups) != 0 {
		t.Logf("User still has claim groups: %v", user.Edges.ClaimGroups)
		t.Fail()
	}

	if user.Username != "ldapuser" || user.Fname != "ldap" ||
			user.Lname != "user" || user.Email != "luser@hppr.dev" ||
			user.Source != "LDAP:main_ldap" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}


// User is in the database, all linked groups were removed, but has a local group assigned
func TestLDAPUpdateUserClaimsRemovedAllGroupsWithRemainingLocal(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("main_ldap", "deleted_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("main_ldap", "deleted_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
		tu.StdLogger(),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	_, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx)
	if err == nil {
		t.Fatal("An error did not occurred getting user claims")
	}
	if errors.Is(usr.AuthenticationError, err) {
		t.Fatal("An authentication error was returned by update user claims, this will stop local look ups.")
	}

	user, claims := getUserAndClaims("ldapuser", ctx)

	if len(claims) != 1 {
		t.Fatalf("Claims length did not match: %v", claims)
	}

	c := claims[0]

	if c.ShortName != "sme" || c.Value != "thg" || c.Name != "some" {
		t.Logf("Claim did not match: %v", c)
		t.Fail()
	}

	if user.Username != "ldapuser" || user.Fname != "ldap" ||
			user.Lname != "user" || user.Email != "luser@hppr.dev" ||
			user.Source != "LDAP:main_ldap" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// Should return an authentication error when an ldap user tries to login, but the connection is unavailable
func TestLDAPUpdateUserClaimsReturnsAnErrorForLDAPUserWhenLDAPConnectionError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		ConnectError(fmt.Errorf("Cannot connect")),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should give claims to local users when ldap is unavailable
func TestLDAPUpdateUserClaimsReturnsAuthSourceError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "local@hppr.dev"),
				tu.Password("imalocal"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("local group", "local group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		ConnectError(fmt.Errorf("Cannot connect")),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	_, err := ldapProvider.UpdateUserClaims("localuser", "imalocal", ctx)
	if err == nil {
		t.Fatalf("An error did not occur: %v", err)
	} else if !errors.Is(err, usr.AuthSourceError) {
		t.Fatalf("Ldap did not return an AuthSourceError: %v", err)
	}
}

// Should return local user claims when the user is not in LDAP
func TestLDAPUpdateUserClaimsReturnsUserNotFoundErrorWhenUserDoesNotExistInLDAP(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "local@hppr.dev"),
				tu.Password("imalocal"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("local group", "local group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	_, err := ldapProvider.UpdateUserClaims("localuser", "imalocal", ctx)
	if err == nil {
		t.Fatalf("An error occured: %v", err)
	} else if !errors.Is(err, usr.UserNotFoundError) {
		t.Fatalf("An unexpected error occurred when expecting UserNotFoundError: %v", err)
	}
}

// Should return error when user exists in ldap, but bad password is given
func TestLDAPUpdateUserClaimsReturnsErrorWithBadLDAPUserPassword(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "whoopsy", ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return error when bad bind user
func TestLDAPUpdateUserClaimsReturnsErrorWithBadLDAPBindUserPassword(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP:main_ldap"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	groupTemplate, userTemplate := createTemplates()
	ldapProvider := usr.NewLDAPUserProvider(
		"some ldap",
		"ldap://someldap.server",
		"adminuser", "adminpass", "", "group_name", "", "first_name", "last_name", "email",
		0,
		groupTemplate, userTemplate,
	)

	ldapProvider.SetConnector(conn)
	ldapProvider.BindUserPassword = "bad"

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when user templates are configured badly
func TestLDAPUpdateUserClaimsReturnsErrorWithMalformedUserFilterTemplate(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
	)

	groupTemplate, userTemplate := createTemplates()
	ldapProvider := usr.NewLDAPUserProvider(
		"some ldap",
		"ldap://someldap.server",
		"adminuser", "adminpass", "", "group_name", "", "first_name", "last_name", "email",
		0,
		groupTemplate, userTemplate,
	)

	ldapProvider.SetConnector(conn)

	ldapProvider.UserFilter.Parse("bad:{{ .NOTEXIST }}")

	if _, err := ldapProvider.UpdateUserClaims("foo", "bar", ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when group templates are configured badly
func TestLDAPUpdateUserClaimsReturnsErrorWithMalformedGroupFilterTemplate(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
	)

	groupTemplate, userTemplate := createTemplates()
	ldapProvider := usr.NewLDAPUserProvider(
		"some ldap",
		"ldap://someldap.server",
		"adminuser", "adminpass", "", "group_name", "", "first_name", "last_name", "email",
		0,
		groupTemplate, userTemplate,
	)

	ldapProvider.SetConnector(conn)

	ldapProvider.GroupFilter.Parse("bad:{{ .NOTEXIST }}")

	if _, err := ldapProvider.UpdateUserClaims("foo", "bar", ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when the user table does not accept a new user
func TestLDAPUpdateUserClaimsReturnsErrorWhenCouldNotCreateUser(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
			),
			tu.ReturnsMutateErrors("user"),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error when trying to create a new ldap user with bad attributes
func TestLDAPUpdateUserClaimsReturnsAnErrorOnMisconfiguredAttributes(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
	)

	groupTemplate, userTemplate := createTemplates()
	ldapProvider := usr.NewLDAPUserProvider(
		"some ldap",
		"ldap://someldap.server",
		"adminuser", "adminpass", "", "group_name", "", "first_name", "last_name", "email",
		0,
		groupTemplate, userTemplate,
	)

	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error if group search fails
func TestLDAPUpdateUserClaimsReturnsAnErrorOnGroupSearchError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
		LDAPGroup("ldapuser", "ldap_group"),
		GroupSearchError(fmt.Errorf("Search failed")),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error user has no LDAP groups
func TestLDAPUpdateUserClaimsReturnsAnErrorOnNoLDAPGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("main_ldap", "ldap_group"),
					tu.Claim(
						tu.ClaimInfo("user", "usr", "yes", "Grants user"),
					),
				),
				tu.Group(
					tu.GroupInfo("local users", "some local group"),
					tu.Claim(
						tu.ClaimInfo("some", "sme", "thg", "Grants something"),
					),
				),
			),
		),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
		// User to login
		LDAPUser("ldapuser", "ldap", "user", "luser@hppr.dev", "luserpass"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

func TestLDAPConnectorReturnsErrorWhenBadURL(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))
	groupTemplate, userTemplate := createTemplates()
	ldapProvider := usr.NewLDAPUserProvider(
		"some ldap",
		"this is a bad url!",
		"", "", "", "", "", "", "", "",
		0,
		groupTemplate, userTemplate,
	)

	if _, err := ldapProvider.UpdateUserClaims("ldapuser", "luserpass", ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

func createLDAPProvider() LDAPUserProvider {
	groupTemplate, userTemplate := createTemplates()

	return usr.NewLDAPUserProvider(
		"main_ldap",
		"ldap://someldap.server",
		"adminuser", "adminpass", "", "group_name", "", "first_name", "last_name", "email",
		0,
		groupTemplate, userTemplate,
	)
}

func createTemplates() (*template.Template, *template.Template) {
	groupTemplate := template.New("group-filter")
	groupTemplate.Parse("groupFilter:{{ .Username }}")

	userTemplate := template.New("user-filter")
	userTemplate.Parse("userFilter:{{ .Username }}")

	return groupTemplate, userTemplate
}

func getUserAndClaims(username string, ctx context.Context) (*ent.User, ent.Claims) {
	user := ent.FromContext(ctx).User.Query().
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
		OnlyX(ctx)
	var allClaims ent.Claims
	for _, group := range user.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	return user, allClaims
}
