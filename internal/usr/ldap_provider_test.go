package usr_test

import (
	"fmt"
	tu "stoke/internal/testutil"
	"stoke/internal/usr"
	"testing"
	"text/template"
)

// User does not exist in the local database yet
// User has claims from linked groups
func TestLDAPGetUserClaimsHappy(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "user@local"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
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

	user, claims, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx)
	if err != nil {
		t.Fatalf("Failed to get user claims: %v", err)
	}

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
func TestLDAPGetUserClaimsNoLinkedGroups(t *testing.T) {
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

	_, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx)
	if err == nil {
		t.Fatalf("Was able to get claims for user with no linked groups")
	}
}

// User is in the database but they have a linked group that is no longer linked to their LDAP groups
func TestLDAPGetUserClaimsRemovedGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("other_ldap_group"),
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

	user, claims, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx)
	if err != nil {
		t.Fatalf("Failed to get user claims: %v", err)
	}

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
			user.Source != "LDAP" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// User is in the database but they have a linked group that is no longer linked to their LDAP groups
// But the database cannot be updated
func TestLDAPGetUserClaimsRemovedGroupsReturnsErrorDatabaseFailure(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("ldap_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("other_ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("Did not return an error")
	}
}

// User is in the database but the database cannot be read
func TestLDAPGetUserClaimsReturnsErrorDatabaseFailure(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("Did not return an error")
	}
}

// User is in the database but all linked groups were removed
func TestLDAPGetUserClaimsRemovedAllGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("deleted_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("deleted_group"),
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

	user, claims, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx)
	if err != nil {
		t.Fatalf("GetUserClaims returned an error: %v", err)
	}

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
			user.Source != "LDAP" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}


// User is in the database, all linked groups were removed, but has a local group assigned
func TestLDAPGetUserClaimsRemovedAllGroupsWithRemainingLocal(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("admin group", "administrator group"),
					tu.LDAPLink("deleted_group"),
					tu.Claim(
						tu.ClaimInfo("admin", "adm", "yes", "Grants admin"),
					),
				),
				tu.Group(
					tu.GroupInfo("other group", "other group"),
					tu.LDAPLink("deleted_group"),
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

	user, claims, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx)
	if err != nil {
		t.Fatalf("An error occurred getting user claims: %v", err)
	}

	if len(claims) != 1 {
		t.Logf("Claims length did not match: %v", claims)
		t.Fail()
	}

	c := claims[0]

	if c.ShortName != "sme" || c.Value != "thg" || c.Name != "some" {
		t.Logf("Claim did not match: %v", c)
		t.Fail()
	}

	if user.Username != "ldapuser" || user.Fname != "ldap" ||
			user.Lname != "user" || user.Email != "luser@hppr.dev" ||
			user.Source != "LDAP" || user.Password != "" || user.Salt != "" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// Should return an authentication error when an ldap user tries to login, but the connection is unavailable
func TestLDAPGetUserClaimsReturnsAnErrorForLDAPUserWhenLDAPConnectionError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should give claims to local users when ldap is unavailable
func TestLDAPGetUserClaimsReturnsClaimsForLocalUserWhenLDAPConnectionError(t *testing.T) {
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

	user, claims, err := ldapProvider.GetUserClaims("localuser", "imalocal", true, ctx)
	if err != nil {
		t.Fatalf("An error occured: %v", err)
	}

	if len(claims) != 1 {
		t.Fatalf("Claims length did not match: %v", claims)
	}

	c := claims[0]
	if c.Name != "user" || c.ShortName != "usr" || c.Value != "yes" {
		t.Logf("Claim did not match: %v", c)
		t.Fail()
	}

	if user.Username != "localuser" || user.Fname != "local" ||
			user.Lname != "user" || user.Email != "local@hppr.dev" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}

}

// Should return local user claims when the user is not in LDAP
func TestLDAPGetUserClaimsReturnsLocalClaimsWhenUserNotFoundInLDAP(t *testing.T) {
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

	user, claims, err := ldapProvider.GetUserClaims("localuser", "imalocal", true, ctx)
	if err != nil {
		t.Fatalf("An error occured: %v", err)
	}

	if len(claims) != 1 {
		t.Fatalf("Claims length did not match: %v", claims)
	}

	c := claims[0]
	if c.Name != "user" || c.ShortName != "usr" || c.Value != "yes" {
		t.Logf("Claim did not match: %v", c)
		t.Fail()
	}

	if user.Username != "localuser" || user.Fname != "local" ||
			user.Lname != "user" || user.Email != "local@hppr.dev" {
		t.Logf("User did not match: %v", user)
		t.Fail()
	}
}

// Should return error when user exists in ldap, but bad password is given
func TestLDAPGetUserClaimsReturnsErrorWithBadLDAPUserPassword(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "whoopsy", true, ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return error when bad bind user
func TestLDAPGetUserClaimsReturnsErrorWithBadLDAPBindUserPassword(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("ldap", "user", "ldapuser", "luser@hppr.dev"),
				tu.Source("LDAP"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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
	ldapProvider.BindUserPassword = "bad"

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when user templates are configured badly
func TestLDAPGetUserClaimsReturnsErrorWithMalformedUserFilterTemplate(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	ldapProvider.UserFilter.Parse("bad:{{ .NOTEXIST }}")

	if _, _, err := ldapProvider.GetUserClaims("foo", "bar", true, ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when group templates are configured badly
func TestLDAPGetUserClaimsReturnsErrorWithMalformedGroupFilterTemplate(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t),
	)

	conn := NewMockLDAPServer(
		// Bind user
		LDAPUser("adminuser", "admin", "user", "admin@hppr.dev", "adminpass"),
	)

	ldapProvider := createLDAPProvider()
	ldapProvider.SetConnector(conn)

	ldapProvider.GroupFilter.Parse("bad:{{ .NOTEXIST }}")

	if _, _, err := ldapProvider.GetUserClaims("foo", "bar", true, ctx); err == nil {
		t.Fatal("An error did not occur")
	}
}

// Should return an error when the user table does not accept a new user
func TestLDAPGetUserClaimsReturnsErrorWhenCouldNotCreateUser(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.StdLogger(),
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error when trying to create a new ldap user with bad attributes
func TestLDAPGetUserClaimsReturnsAnErrorOnMisconfiguredAttributes(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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
	ldapProvider.FirstNameField = "cantfindme"

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error if group search fails
func TestLDAPGetUserClaimsReturnsAnErrorOnGroupSearchError(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

// Should return an error user has no LDAP groups
func TestLDAPGetUserClaimsReturnsAnErrorOnNoLDAPGroups(t *testing.T) {
	ctx := tu.NewMockContext(
		tu.WithDatabase(t,
			tu.User(
				tu.UserInfo("local", "user", "localuser", "luser@hppr.dev"),
				tu.Source("LOCAL"),
				tu.Group(
					tu.GroupInfo("ldap user group", "ldap group"),
					tu.LDAPLink("ldap_group"),
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

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

func TestLDAPConnectorReturnsErrorWhenBadURL(t *testing.T) {
	ctx := tu.NewMockContext(tu.WithDatabase(t))
	ldapProvider := createLDAPProvider()

	ldapProvider.ServerURL = "this is a bad url!"
	ldapProvider.DefaultConnector()

	if _, _, err := ldapProvider.GetUserClaims("ldapuser", "luserpass", true, ctx); err == nil {
		t.Fatal("An error did not occur.")
	}
}

func createLDAPProvider() usr.LDAPUserProvider {
	groupTemplate, userTemplate := createTemplates()

	return usr.LDAPUserProvider{
		BindUserDN:       "adminuser",
		BindUserPassword: "adminpass",
		GroupFilter:      groupTemplate,
		GroupAttribute:   "group_name",
		UserFilter:       userTemplate,
		FirstNameField:   "first_name",
		LastNameField:    "last_name",
		EmailField:       "email",
		LocalProvider:    usr.LocalProvider{},
	}
}

func createTemplates() (*template.Template, *template.Template) {
	groupTemplate := template.New("group-filter")
	groupTemplate.Parse("groupFilter:{{ .Username }}")

	userTemplate := template.New("user-filter")
	userTemplate.Parse("userFilter:{{ .Username }}")

	return groupTemplate, userTemplate
}
