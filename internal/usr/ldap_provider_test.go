package usr_test

import (
	tu "stoke/internal/testutil"
	"testing"
)

func TestLDAPGetUserClaimsHappy(t *testing.T) {
	_ = tu.NewMockContext()
	// (user, password string, ctx context.Context) (*ent.User, ent.Claims, error)
}
