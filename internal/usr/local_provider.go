package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/user"

	"golang.org/x/crypto/argon2"
)

type LocalProvider struct {
	DB     *ent.Client
}

func (l LocalProvider) Init() error {
	return l.checkForSuperUser()
}

func (l LocalProvider) AddUser(fname, lname, email, username, password string, superUser bool) error {
	salt := l.genSalt()
	userInfo, err := l.DB.User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSalt(salt).
		SetPassword(l.hashPass(password, salt)).
		Save(context.Background())
	if err != nil {
		return nil
	}
	
	_, err = l.DB.ClaimGroup.Create().
		AddUsers(userInfo).
		SetName(username).
		SetDescription(fmt.Sprintf("%s's group", username)).
		SetIsUserGroup(true).
		Save(context.Background())
	if err != nil {
		return nil
	}

	if superUser {
		superGroup, err := l.getOrCreateSuperGroup()
		if err != nil {
			return err
		}
		_, err = superGroup.Update().AddUsers(userInfo).Save(context.Background())
		return err
	}

	return nil
}

func (l LocalProvider) GetUserClaims(username, password string) (ent.Claims, error) {
	usr, err := l.DB.User.Query().
		Where(
			user.Or(
				user.UsernameEQ(username),
				user.EmailEQ(username),
			),
		).
		WithClaimGroups(func (q *ent.ClaimGroupQuery) {
			q.WithClaims()
		}).
		Only(context.Background())
	if err != nil {
		return nil, err
	}
	if l.hashPass(password, usr.Salt) != usr.Password {
		return nil, fmt.Errorf("Bad Password")
	}
	var allClaims ent.Claims
	for _, group := range usr.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	return allClaims, nil
}

func (l LocalProvider) hashPass(pass, salt string) string {
		return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func (l LocalProvider) genSalt() string {
	saltBytes := make([]byte, 32)
	rand.Read(saltBytes)
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func (l LocalProvider) getOrCreateSuperGroup() (*ent.ClaimGroup, error) {
	superGroup, err := l.DB.ClaimGroup.Query().
		Where(
			claimgroup.HasClaimsWith(
				claim.And(
					claim.ShortNameEQ("srol"),
					claim.ValueEQ("spr"),
				),
			),
		).
		WithUsers().
		Only(context.Background())

	if ent.IsNotFound(err) {
		log.Println("Stoke superusers not found. Creating...")
		superClaim := l.DB.Claim.Create().
			SetName("Stoke Super User").
			SetDescription("Grants superuser management access to the stoke server").
			SetShortName("srol").
			SetValue("spr").
			SaveX(context.Background())

		superGroup = l.DB.ClaimGroup.Create().
			AddClaims(superClaim).
			SetName("Stoke Superusers").
			SetDescription("Stoke server superusers").
			SaveX(context.Background())
	} else if err != nil {
		return nil, err
	}
	return superGroup, nil
}

func (l LocalProvider) checkForSuperUser() error {
	superGroup, err := l.getOrCreateSuperGroup()
	if err != nil {
		return err
	}
	if len(superGroup.Edges.Users) == 0 {
		randomPass := l.genSalt()
		l.AddUser("Stoke", "Admin", "sadmin@localhost", "sadmin", randomPass, true)
		log.Printf("Created superuser 'sadmin' with password %s", randomPass)
	}
	return nil
}
