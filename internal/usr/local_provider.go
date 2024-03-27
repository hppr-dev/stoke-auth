package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"stoke/internal/ent"
	"stoke/internal/ent/user"

	"golang.org/x/crypto/argon2"
)

type LocalProvider struct {
	Schema string
	DB     *ent.Client
}

func (l LocalProvider) Init() error {
	return nil
}

func (l LocalProvider) AddUser(fname, lname, email, username, pass string) error {
	salt := l.genSalt()
	_, err := l.DB.User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSalt(salt).
		SetPassword(l.hashPass(pass, salt)).
		Save(context.Background())
	return err
}

func (l LocalProvider) ValidateUser(username, pass string) bool {
	user, err := l.DB.User.Query().
		Where(
			user.Or(
				user.UsernameEQ(username),
				user.EmailEQ(username),
			),
		).
		Only(context.Background())
	if err != nil {
		return false
	}
	return user.Password == l.hashPass(pass, user.Salt)
}

func (l LocalProvider) hashPass(pass, salt string) string {
		return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func (l LocalProvider) genSalt() string {
	saltBytes := make([]byte, 32)
	rand.Read(saltBytes)
	return base64.StdEncoding.EncodeToString(saltBytes)
}
