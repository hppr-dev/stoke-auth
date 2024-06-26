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

func HashPass(pass, salt string) string {
	return base64.StdEncoding.EncodeToString(argon2.IDKey([]byte(pass), []byte(salt), 2, 19*1024, 1, 64))
}

func GenSalt() string {
	saltBytes := make([]byte, 32)
	rand.Read(saltBytes)
	return base64.StdEncoding.EncodeToString(saltBytes)
}

func ProviderFromCtx(ctx context.Context) Provider {
	return ctx.Value("user-provider").(Provider)
}

func (l LocalProvider) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, "user-provider", l)
}

func (l LocalProvider) AddUser(fname, lname, email, username, password string, _ bool, ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	ctx, span := tel.GetTracer().Start(ctx, "LoginApiHandler.ServeHTTP")
	defer span.End()

	logger.Info().
		Func(otelzerolog.AddTracingContext(span)).
		Str("fname", fname).
		Str("lname", lname).
		Str("username", username).
		Str("email", email).
		Msg("Creating user")

	salt := GenSalt()
	_, err := ent.FromContext(ctx).User.Create().
		SetFname(fname).
		SetLname(lname).
		SetEmail(email).
		SetUsername(username).
		SetSource("LOCAL").
		SetSalt(salt).
		SetPassword(HashPass(password, salt)).
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
	
	return nil
}

func (l LocalProvider) GetUserClaims(username, password string, verify bool, ctx context.Context) (*ent.User, ent.Claims, error) {
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

	if verify && HashPass(password, usr.Salt) != usr.Password {
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
		Interface("claims", allClaims).
		Msg("Claims found")
	return usr, allClaims, nil
}

func (l LocalProvider) getOrCreateSuperGroup(ctx context.Context) (*ent.ClaimGroup, error) {
	logger := zerolog.Ctx(ctx)

	client := ent.FromContext(ctx)

	superGroup, err := client.ClaimGroup.Query().
		Where(
			claimgroup.HasClaimsWith(
				claim.And(
					claim.ShortNameEQ("stk"),
					claim.ValueEQ("S"),
				),
			),
		).
		WithUsers().
		First(ctx)

	if ent.IsNotFound(err) {
		logger.Info().
			Msg("Stoke superusers not found. Creating...")

		superClaim, err := client.Claim.Create().
			SetName("Stoke Super User").
			SetDescription("Grants superuser management access to the stoke server").
			SetShortName("stk").
			SetValue("S").
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

func (l LocalProvider) CheckCreateForSuperUser(ctx context.Context) error {
	superGroup, err := l.getOrCreateSuperGroup(ctx)
	if err != nil {
		return err
	}
	if len(superGroup.Edges.Users) == 0 {
		randomPass := GenSalt()
		l.AddUser("Stoke", "Admin", "sadmin@localhost", "sadmin", randomPass, true, ctx)
		ent.FromContext(ctx).User.Update().
			Where(user.UsernameEQ("sadmin")).
			AddClaimGroups(superGroup).
			SaveX(ctx)
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

	if !force && HashPass(oldPassword, usr.Salt) != usr.Password {
		logger.Debug().
			Func(otelzerolog.AddTracingContext(span)).
			Str("username", username).
			Msg("User old password did not match")
		return fmt.Errorf("Bad Password")
	}

	newSalt := GenSalt()

	_, err = usr.Update().
		SetSalt(newSalt).
		SetPassword(HashPass(newPassword, newSalt)).
		Save(ctx)
	return err
}

func (l LocalProvider) CheckCreateForStokeClaims(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	if err := l.checkCreateClaim("Read Claims", "Grants read access to claims", "c", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read claims claim")
	}
	if err := l.checkCreateClaim("Write Claims", "Grants read/write access to claims", "C", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read/write claims claim")
	}
	if err := l.checkCreateClaim("Read Users", "Grants read access to users", "u", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read users claim")
	}
	if err := l.checkCreateClaim("Write Users", "Grants read/write access to users", "U", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read/write users claim")
	}
	if err := l.checkCreateClaim("Read Groups", "Grants read access to groups", "g", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read groups claim")
	}
	if err := l.checkCreateClaim("Write Groups", "Grants read/write access to groups", "G", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create read/write groups claim")
	}
	if err := l.checkCreateClaim("Super Read", "Grants read access to all stoke admin", "s", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create super read claim")
	}
	if err := l.checkCreateClaim("Monitoring Access", "Grants read access to stoke monitoring", "m", ctx); err != nil {
		logger.Warn().Err(err).Msg("Could not create monitoring claim")
	}
	return nil
}

func (l LocalProvider) checkCreateClaim(name, desc, value string, ctx context.Context) error {
	client := ent.FromContext(ctx)

	_, err := client.Claim.Query().
		Where(
			claim.And(
				claim.ShortNameEQ("stk"),
				claim.ValueEQ(value),
			),
		).FirstID(ctx)
		if ent.IsNotFound(err) {
			_, err := client.Claim.Create().
				SetName(name).
				SetDescription(desc).
				SetShortName("stk").
				SetValue(value).
				Save(ctx)
			return err
		}
		return nil
}
