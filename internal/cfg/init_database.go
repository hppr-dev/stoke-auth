package cfg

import (
	"context"
	"fmt"
	"os"
	"stoke/internal/ent"
	"stoke/internal/ent/claim"
	"stoke/internal/ent/claimgroup"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/rs/zerolog"
)

func InitializeDatabaseFromFile(filename string, ctx context.Context) error {
	logger := zerolog.Ctx(ctx).With().Str("filename", filename).Logger()

	content, err := os.ReadFile(filename)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not read database initialization file.")
		return err
	}

	logger.Info().Msg("Initializing database")

	dbInitFile := &databaseInitFile{}
	if err = yaml.Unmarshal(content, dbInitFile); err != nil {
		logger.Error().
			Err(err).
			Msg("Could not unmarshal database initialization file.")
	}

	if err := dbInitFile.validate(); err != nil {
		logger.Error().
			Err(err).
			Msg("Could not start database transaction")
		return err
	}

	tx, err := ent.FromContext(ctx).Tx(ctx)
	if err != nil {
		logger.Error().
			Err(err).
			Msg("Could not start database transaction")
		return err
	}

	for _, c := range dbInitFile.Claims {
		if err := c.writeToDB(tx, ctx); err != nil {
			logger.Error().
				Err(err).
				AnErr("rollbackErr", tx.Rollback()).
				Msg("Failed to initialize claims in database")
			return err
		}
	}

	for _, g := range dbInitFile.Groups {
		if err := g.writeToDB(tx, ctx); err != nil {
			logger.Error().
				Err(err).
				AnErr("rollbackErr", tx.Rollback()).
				Msg("Failed to initialize claims in database")
			return err
		}
	}

	for _, u := range dbInitFile.Users {
		if err := u.writeToDB(tx, ctx); err != nil {
			logger.Error().
				Err(err).
				AnErr("rollbackErr", tx.Rollback()).
				Msg("Failed to initialize claims in database")
			return err
		}
	}

	return tx.Commit()
}

type databaseInitFile struct {
	Claims []initClaim `json:"claims"`
	Groups []initGroup `json:"groups"`
	Users []initUser   `json:"users"`
}

func (d *databaseInitFile) validate() error {
	claimMap := make(map[string]bool)
	groupMap := make(map[string]bool)

	for _, claim := range d.Claims {
		claimMap[claim.Name] = true
	}

	for _, group := range d.Groups {
		groupMap[group.Name] = true
		for _, claimName := range group.Claims {
			if _, ok := claimMap[claimName]; !ok {
				return fmt.Errorf("Could not find claim %s for group %s", claimName, group.Name)
			}
		}
	}

	for _, user := range d.Users {
		for _, groupName := range user.Groups {
			if _, ok := groupMap[groupName]; !ok {
				return fmt.Errorf("Could not find group %s for user %s", groupName, user.Username)
			}
		}
	}

	return nil
}

type initUser struct {
	FName string         `json:"first_name"`
	LName string         `json:"last_name"`
	Email string         `json:"email"`
	Username string      `json:"username"`
	PasswordHash string  `json:"password_hash"`
	PasswordSalt string  `json:"password_salt"`
	Groups []string      `json:"groups"`
}

func (u initUser) writeToDB(tx *ent.Tx, ctx context.Context) error {
	groups, err :=  tx.ClaimGroup.Query().
		Where(
			claimgroup.NameIn(u.Groups...),
		).
		IDs(ctx)
	if err != nil {
		return err
	}

	_, err = tx.User.Create().
		SetFname(u.FName).
		SetLname(u.LName).
		SetEmail(u.Email).
		SetUsername(u.Username).
		SetPassword(u.PasswordHash).
		SetSalt(u.PasswordSalt).
		SetSource("LOCAL").
		AddClaimGroupIDs(groups...).
		Save(ctx)
	return err
}

type initGroup struct {
	Name string        `json:"name"`
	Description string `json:"description"`
	Claims []string    `json:"claims"`
	Links []string     `json:"links"`
}

func (g initGroup) writeToDB(tx *ent.Tx, ctx context.Context) error {
	claims, err := tx.Claim.Query().
		Where(
			claim.NameIn(g.Claims...),
		).
		IDs(ctx)
	if err != nil {
		return err
	}

	group, err := tx.ClaimGroup.Create().
		SetName(g.Name).
		SetDescription(g.Description).
		AddClaimIDs(claims...).
		Save(ctx)
	if err != nil {
		return err
	}

	if len(g.Links) > 0 {
		for _, link := range g.Links {
			_, err = tx.GroupLink.Create().
				SetClaimGroup(group).
				SetType("LDAP").
				SetResourceSpec(link).
				Save(ctx)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type initClaim struct {
	Name string        `json:"name"`
	Description string `json:"description"`
	ShortName string   `json:"short_name"`
	Value string       `json:"value"`
	ClaimString string `json:"claim_string"`
}

func (c initClaim) writeToDB(tx *ent.Tx, ctx context.Context) error {
	if c.ClaimString != "" {
		splitStr := strings.Split(c.ClaimString, "=")
		if len(splitStr) != 2 {
			return fmt.Errorf("Bad claim_string format. Got %s but expected value like short_name=value.", c.ClaimString)
		}
		c.ShortName = splitStr[0]
		c.Value = splitStr[1]
	}

	_, err := tx.Claim.Create().
		SetDescription(c.Description).
		SetName(c.Name).
		SetValue(c.Value).
		SetShortName(c.ShortName).
		Save(ctx)
	return err
}
