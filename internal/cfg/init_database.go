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
			Msg("Could not start validate db init file")
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
				Interface("claim", c).
				Msg("Failed to initialize claims in database")
			return err
		}
	}

	for _, g := range dbInitFile.Groups {
		if err := g.writeToDB(tx, ctx); err != nil {
			logger.Error().
				Err(err).
				AnErr("rollbackErr", tx.Rollback()).
				Interface("group", g).
				Msg("Failed to initialize claims in database")
			return err
		}
	}

	for _, u := range dbInitFile.Users {
		if err := u.writeToDB(tx, ctx); err != nil {
			logger.Error().
				Err(err).
				AnErr("rollbackErr", tx.Rollback()).
				Interface("user", u).
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

	for i := range d.Claims {
		if err := d.Claims[i].parseShortHand(); err != nil {
			return err
		}

		claimMap[d.Claims[i].Name] = true
	}

	for i := range d.Groups {
		if err := d.Groups[i].parseShortHand(); err != nil {
			return err
		}

		groupMap[d.Groups[i].Name] = true
		for _, claimName := range d.Groups[i].Claims {
			if _, ok := claimMap[claimName]; !ok {
				return fmt.Errorf("Could not find claim %s for group %s", claimName, d.Groups[i].Name)
			}
		}
	}

	for i := range d.Users {
		if err := d.Users[i].parseShortHand(); err != nil {
			return err
		}
		for _, groupName := range d.Users[i].Groups {
			if _, ok := groupMap[groupName]; !ok {
				return fmt.Errorf("Could not find group %s for user %s", groupName, d.Users[i].Username)
			}
		}
	}

	return nil
}

type initUser struct {
	Username string      `json:"username"`
	FName string         `json:"first_name"`
	LName string         `json:"last_name"`
	Email string         `json:"email"`
	PasswordHash string  `json:"password_hash"`
	PasswordSalt string  `json:"password_salt"`
	Groups []string      `json:"groups"`

	UserString string    `json:"user_string"`
}

func (u *initUser) parseShortHand() error {
	if u.UserString != "" {
		splitStr := strings.Split(u.UserString, ",")
		if len(splitStr) != 4 {
			return fmt.Errorf("Bad user_string format. Got %s but expected value like username,first_name,last_name,email", u.UserString)
		}
		u.Username = strings.Trim(splitStr[0], " ")
		u.FName = strings.Trim(splitStr[1], " ")
		u.LName = strings.Trim(splitStr[2], " ")
		u.Email = strings.Trim(splitStr[3], " ")
	}

	return nil
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

	GroupString string `json:"group_string"`
}

func (g *initGroup) parseShortHand() error {
	if g.GroupString != "" {
		splitStr := strings.Split(g.GroupString, ",")
		if len(splitStr) != 2 {
			return fmt.Errorf("Bad claim_string format. Got %s but expected value like name,description", g.GroupString)
		}
		g.Name = strings.Trim(splitStr[0], " ")
		g.Description = strings.Trim(splitStr[1], " ")
	}
	return nil
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

func (c *initClaim) parseShortHand() error {
	if c.ClaimString != "" {
		splitStr := strings.Split(c.ClaimString, ",")
		if len(splitStr) != 4 {
			return fmt.Errorf("Bad claim_string format. Got %s but expected value like name,description,short_name,value", c.ClaimString)
		}
		c.Name = strings.Trim(splitStr[0], " ")
		c.Description = strings.Trim(splitStr[1], " ")
		c.ShortName = strings.Trim(splitStr[2], " ")
		c.Value = strings.Trim(splitStr[3], " ")
	}
	return nil
}

func (c initClaim) writeToDB(tx *ent.Tx, ctx context.Context) error {

	_, err := tx.Claim.Create().
		SetDescription(c.Description).
		SetName(c.Name).
		SetValue(c.Value).
		SetShortName(c.ShortName).
		Save(ctx)
	return err
}
