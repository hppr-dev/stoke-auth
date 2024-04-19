package usr

import (
	"context"
	"stoke/internal/ent"
)

type Provider interface {
	Init(context.Context) error
	GetUserClaims(user, password string, ctx context.Context) (*ent.User, ent.Claims, error)
  AddUser(fname, lname, email, username, password string, superuser bool, ctx context.Context) error
  UpdateUserPassword(username, oldPassword, newPassword string, force bool, ctx context.Context) error
}
