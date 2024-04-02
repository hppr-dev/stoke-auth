package usr

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/user"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/argon2"
)

type LocalProvider struct {
	DB     *ent.Client
}

func (l LocalProvider) Init() error {
	return l.checkForSuperUser()
}

func (l LocalProvider) AddUser(provider ProviderType, fname, lname, email, username, password string, superUser bool) error {
	if provider != LOCAL {
		return ProviderTypeNotSupported
	}

	logger.Info().
			Str("fname", fname).
			Str("lname", lname).
			Str("username", username).
			Str("email", email).
			Bool("superuser", superUser).
			Msg("Creating user")

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
		logger.Error().
			Err(err).
			Str("username", username).
			Str("email", email).
			Msg("Could not create user")
		return err
	}
	
	_, err = l.DB.ClaimGroup.Create().
		AddUsers(userInfo).
		SetName(username).
		SetDescription(fmt.Sprintf("%s's group", username)).
		SetIsUserGroup(true).
		Save(context.Background())
	if err != nil {
		logger.Error().
			Err(err).
			Str("username", username).
			Str("email", email).
			Msg("Could not create user claim group")
		return err
	}

	if superUser {
		superGroup, err := l.getOrCreateSuperGroup()
		if err != nil {
		logger.Error().
			Err(err).
			Str("username", username).
			Str("email", email).
			Msg("Could not get superuser group")
			return err
		}

		_, err = superGroup.Update().AddUsers(userInfo).Save(context.Background())
		if err != nil {
			logger.Error().
				Err(err).
				Str("username", username).
				Str("email", email).
				Msg("Could add user to super group")
				return err
		}
	}

	return nil
}

func (l LocalProvider) GetUserClaims(username, password string) (ent.Claims, error) {
	logger.Debug().
		Str("username", username).
		Msg("Getting user claims")

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
		logger.Error().
			Err(err).
			Str("username", username).
			Msg("Could not find user")
		return nil, err
	}

	if l.hashPass(password, usr.Salt) != usr.Password {
		logger.Debug().Str("username", username).Msg("User password did not match")
		return nil, fmt.Errorf("Bad Password")
	}

	var allClaims ent.Claims
	for _, group := range usr.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	logger.Debug().
		Str("username", username).
		Func(func (e *zerolog.Event) {
			var values []string
			for _, c := range allClaims {
				values = append(values, c.ShortName + ":" + c.Value)
			}
			e.Strs("claims", values)
		}).
		Msg("Claims found")
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
		logger.Info().Msg("Stoke superusers not found. Creating...")
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
		l.AddUser(LOCAL, "Stoke", "Admin", "sadmin@localhost", "sadmin", randomPass, true)
		logger.Info().
			Str("password", randomPass).
			Msg("Created superuser 'sadmin'")
	}
	return nil
}

func (l LocalProvider) UpdateUser(provider ProviderType, fname, lname, email, username, password string) error {
	// TODO
	return fmt.Errorf("NOT IMPLEMENTED")
}
