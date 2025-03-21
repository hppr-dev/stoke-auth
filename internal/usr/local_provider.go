package usr

import (
	"context"
	"fmt"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/schema/policy"
	"stoke/internal/ent/user"
	"stoke/internal/tel"

	"github.com/rs/zerolog"
	"github.com/vincentfree/opentelemetry/otelzerolog"
)

type localProvider struct {}

func (l *localProvider) AddUser(fname, lname, email, username, password string, ctx context.Context) error {
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
		SetSource(LOCAL_SOURCE).
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

func (l *localProvider) GetUserClaims(username, password string, u *ent.User, ctx context.Context) (*ent.User, ent.Claims, error) {
	ctx, span := tel.GetTracer().Start(ctx, "LocalUserProvider.GetUserClaims")
	defer span.End()

	logger := zerolog.Ctx(ctx).With().
		Str("component", "LocalUserProvider.GetUserClaims").
		Logger()

	var allClaims ent.Claims
	var err error

	logger.Debug().Interface("user", u).Msg("Getting local claims")

	// Find the user if we haven't got one yet
	if u == nil {
		logger.Info().Msg("Finding local user")
		u, err = retreiveLocalUser(username, ctx)
		if err != nil {
			logger.Error().
				Func(otelzerolog.AddTracingContext(span)).
				Err(err).
				Msg("Could not retrieve local claims")
			return nil, nil, err
		}

		// We only need to verify the stored password hash if the user is local.
		// The password will be blank for all other sources.
		// Other provider sources must verify the password before retreiving claims from local
		if HashPass(password, u.Salt) != u.Password {
			logger.Debug().
				Func(otelzerolog.AddTracingContext(span)).
				Msg("User password did not match")
			return nil, nil, fmt.Errorf("Bad Password")
		}
	}

	allClaims = allUserClaims(u)
	logger.Debug().
		Func(otelzerolog.AddTracingContext(span)).
		Str("username", u.Username).
		Interface("claims", allClaims).
		Msg("Claims found")
	return u, allClaims, nil
}

func (l *localProvider) getOrCreateSuperGroup(ctx context.Context) (*ent.ClaimGroup, error) {
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

func (l *localProvider) CheckCreateForSuperUser(ctx context.Context) error {
	ctx = policy.BypassDatabasePolicies(ctx)
	superGroup, err := l.getOrCreateSuperGroup(ctx)
	logger := zerolog.Ctx(ctx).With().
		Str("component", "CheckCreateForSuperUser").
		Logger()

	if err != nil {
		return err
	}

	if len(superGroup.Edges.Users) == 0 {
		randomPass := GenSalt()
		if err := l.AddUser("Stoke", "Admin", "sadmin@localhost", "sadmin", randomPass, ctx); err != nil {
			logger.Error().
				Err(err).
				Msg("Error while creating user")
		}

		ent.FromContext(ctx).User.Update().
			Where(user.UsernameEQ("sadmin")).
			AddClaimGroups(superGroup).
			SaveX(ctx)
		logger.Info().
			Str("password", randomPass).
			Msg("Created superuser 'sadmin'")
	}
	return nil
}

func (l *localProvider) UpdateUserPassword(username, oldPassword, newPassword string, force bool, ctx context.Context) error {
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

func (l *localProvider) CheckCreateForStokeClaims(ctx context.Context) error {
	logger := zerolog.Ctx(ctx)
	ctx = policy.BypassDatabasePolicies(ctx)
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

func (l *localProvider) checkCreateClaim(name, desc, value string, ctx context.Context) error {
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
