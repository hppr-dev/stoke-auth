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
	"stoke/internal/tel"

	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
	"golang.org/x/crypto/argon2"
)

type LocalProvider struct {}

func (l LocalProvider) Init(ctx context.Context) error {
	return l.checkForSuperUser(ctx)
}

func (l LocalProvider) AddUser(fname, lname, email, username, password string, superUser bool, ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LoginApiHandler.ServeHTTP")
	defer span.End()


	logger.Info().
		Func(otelzerolog.AddTracingContext(span)).
		Str("fname", fname).
		Str("lname", lname).
		Str("username", username).
		Str("email", email).
		Bool("superuser", superUser).
		Msg("Creating user")

	salt := l.genSalt()
	userInfo, err := ent.FromContext(ctx).User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSource("LOCAL").
		SetSalt(salt).
		SetPassword(l.hashPass(password, salt)).
		Save(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Str("email", email).
			Msg("Could not create user")
		return err
	}
	
	if superUser {
		superGroup, err := l.getOrCreateSuperGroup(ctx)
		if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Str("email", email).
			Msg("Could not get superuser group")
			return err
		}

		_, err = superGroup.Update().AddUsers(userInfo).Save(ctx)
		if err != nil {
			logger.Error().
				Func(otelzerolog.AddTracingContext(span)).
				Err(err).
				Str("username", username).
				Str("email", email).
				Msg("Could add user to super group")
				return err
		}
	}

	return nil
}

func (l LocalProvider) GetUserClaims(username, password string, ctx context.Context) (*ent.User, ent.Claims, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LocalUserProvider.GetUserClaims")
	defer span.End()

	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", username).
		Msg("Getting user claims")

	usr, err := ent.FromContext(ctx).User.Query().
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
		Only(ctx)
	if err != nil {
		logger.Error().
			Func(otelzerolog.AddTracingContext(span)).
			Err(err).
			Str("username", username).
			Msg("Could not find user")
		return nil, nil, err
	}

	if usr.Source != "LDAP" && l.hashPass(password, usr.Salt) != usr.Password {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("username", username).
			Msg("User password did not match")
		return nil, nil, fmt.Errorf("Bad Password")
	}

	var allClaims ent.Claims
	for _, group := range usr.Edges.ClaimGroups {
		allClaims = append(allClaims, group.Edges.Claims...)
	}
	logger.Debug().
		Str("username", username).
		Func(otelzerolog.AddTracingContext(span)).
		Func(func (e *zerolog.Event) {
			var values []string
			for _, c := range allClaims {
				values = append(values, c.ShortName + ":" + c.Value)
			}
			e.Strs("claims", values)
		}).
		Msg("Claims found")
	return usr, allClaims, nil
}

func (l LocalProvider) hashPass(pass, salt string) string {
		return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func (l LocalProvider) genSalt() string {
	saltBytes := make([]byte, 32)
	rand.Read(saltBytes)
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func (l LocalProvider) getOrCreateSuperGroup(ctx context.Context) (*ent.ClaimGroup, error) {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LocalUserProvider.getOrCreateSuperGroup")
	defer span.End()

	client := ent.FromContext(ctx)

	superGroup, err := client.ClaimGroup.Query().
		Where(
			claimgroup.HasClaimsWith(
				claim.And(
					claim.ShortNameEQ("srol"),
					claim.ValueEQ("spr"),
				),
			),
		).
		WithUsers().
		First(ctx)

	if ent.IsNotFound(err) {
		logger.Info().
			Func(otelzerolog.AddTracingContext(span)).
			Msg("Stoke superusers not found. Creating...")

		superClaim, err := client.Claim.Create().
			SetName("Stoke Super User").
			SetDescription("Grants superuser management access to the stoke server").
			SetShortName("srol").
			SetValue("spr").
			Save(ctx)
		if err != nil {
			return nil, err
		}

		superGroup, err = client.ClaimGroup.Create().
			AddClaims(superClaim).
			SetName("Stoke Superusers").
			SetDescription("Stoke server superusers").
			Save(ctx)
		if err != nil {
			return nil, err
		}

	} else if err != nil {
		return nil, err
	}
	return superGroup, nil
}

func (l LocalProvider) checkForSuperUser(ctx context.Context) error {
	ctx, span := tel.GetTracer().Start(ctx, "LocalUserProvider.checkForSuperUser")
	defer span.End()

	superGroup, err := l.getOrCreateSuperGroup(ctx)
	if err != nil {
		return err
	}
	if len(superGroup.Edges.Users) == 0 {
		randomPass := l.genSalt()
		l.AddUser("Stoke", "Admin", "sadmin@localhost", "sadmin", randomPass, true, ctx)
		zerolog.Ctx(ctx).Info().
			Str("password", randomPass).
			Msg("Created superuser 'sadmin'")
	}
	return nil
}

func (l LocalProvider) UpdateUserPassword(username, oldPassword, newPassword string, force bool, ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LocalUserProvider.UpdateUser")
	defer span.End()

	usr, err := ent.FromContext(ctx).User.Query().
		Where(user.UsernameEQ(username)).
		Only(context.Background())
	if err != nil {
		return err
	}

	if !force && l.hashPass(oldPassword, usr.Salt) != usr.Password {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("username", username).
			Bool("force", force).
			Msg("User old password did not match")
		return fmt.Errorf("Bad Password")
	}

	newSalt := l.genSalt()
	newPassHash := l.hashPass(newPassword, newSalt)

	_, err = usr.Update().
		SetSalt(newSalt).
		SetPassword(newPassHash).
		Save(ctx)
	return err
}
