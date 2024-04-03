package web

import (
	"errors"
	"fmt"
	"net/http"
	"stoke/internal/cfg"
	"stoke/internal/ctx"
	"stoke/internal/ent"
	"stoke/internal/key"
	"time"

	"github.com/go-faster/jx"
	"github.com/golang-jwt/jwt/v5"
)

type LoginApiHandler struct {
	Context *ctx.Context
}

// Request takes username, password and optionally required_claims.
// required_claims is an object specifying which claim the user must have to receive a token
// If required_claims is not included, a token is granted if the username and password are correct
func (l LoginApiHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		MethodNotAllowed.Write(res)
		return
	}

	var username, password string
	requiredClaims := make(map[string]string)
	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "username":
			username, err = d.Str()
		case "password":
			password, err = d.Str()
		case "required_claims":
			err = d.Obj(func ( d *jx.Decoder, key string) error {
				val, err := d.Str()
				requiredClaims[key] = val
				return err
			})
		default:
			return errors.New("Bad Request")
		}
		return err
	})

	if err != nil || username == "" || password == "" {
		logger.Debug().Err(err).Str("username", username).Msg("Missing body parameters")
		BadRequest.Write(res)
		return
	}
	user, claims, err := l.Context.UserProvider.GetUserClaims(username, password)
	if err != nil {
		logger.Debug().Err(err).Msg("Failed to get claims from provider")
		Unauthorized.Write(res)
		return
	}
	
	claimMap := make(map[string]string)

	for _, claim := range claims {
		claimMap[claim.ShortName] = claim.Value
	}

	for reqKey, reqValue := range requiredClaims {
		userValue, ok := claimMap[reqKey]
		if !ok || userValue != reqValue {
			logger.Debug().
				Str("claimShortName", reqKey).
				Str("requiredValue", reqValue).
				Str("actualValue", userValue).
				Msg("User did not have required claims.")
			Unauthorized.Write(res)
			return
		}
	}

	populateUserInfo(l.Context.Config, user, claimMap)

	token, refresh, err := l.Context.Issuer.IssueToken(key.Claims{
		StokeClaims : claimMap,
		RegisteredClaims: createRegisteredClaims(l.Context.Config),
	})
	if err != nil {
		InternalServerError.Write(res)
		return
	}

	res.Write([]byte(fmt.Sprintf("{\"token\":\"%s\",\"refresh\":\"%s\"}", token, refresh)))
}

func createRegisteredClaims(c cfg.Config) jwt.RegisteredClaims {
	now := time.Now()
	minClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(now.Add(c.Tokens.TokenDuration)),

		// Below fields are omitted if they are not included in config
		Issuer:    c.Tokens.Issuer,
		Subject:   c.Tokens.Subject,
		Audience:  c.Tokens.Audience,
	}

	if c.Tokens.IncludeNotBefore {
		minClaims.NotBefore = jwt.NewNumericDate(now)
	}

	if c.Tokens.IncludeIssuedAt {
		minClaims.IssuedAt = jwt.NewNumericDate(now)
	}

	return minClaims
}

func populateUserInfo(c cfg.Config, user *ent.User, t map[string]string) {
	usernameKey, ok := c.Tokens.UserInfo["username"]
	if ok {
		t[usernameKey] = user.Username
	}

	fnameKey, ok := c.Tokens.UserInfo["first_name"]
	if ok {
		t[fnameKey] = user.Fname
	}

	lnameKey, ok := c.Tokens.UserInfo["last_name"]
	if ok {
		t[lnameKey] = user.Lname
	}

	nameKey, ok := c.Tokens.UserInfo["full_name"]
	if ok {
		t[nameKey] = fmt.Sprintf("%s %s", user.Fname, user.Lname)
	}

	emailKey, ok := c.Tokens.UserInfo["email"]
	if ok {
		t[emailKey] = user.Email
	}
}
