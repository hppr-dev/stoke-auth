package web

import (
	"errors"
	"net/http"
	"stoke/internal/ctx"
	"stoke/internal/usr"

	"github.com/go-faster/jx"
)

type UserHandler struct {
	Context *ctx.Context
}

func (h UserHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
		case http.MethodPost:
			h.handleCreate(res, req)
		case http.MethodPatch:
			h.handleUpdate(res, req)
		default:
			MethodNotAllowed.Write(res)
	}
}

// Takes json with the following fields: fname, lname, email, username, password, provider and superuser
// provider must be either local or ldap
func (h UserHandler) handleCreate(res http.ResponseWriter, req *http.Request) {
	var fname, lname, email, username, password, provider, superuser string

	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "fname":
			fname, err = d.Str()
		case "lname":
			lname, err = d.Str()
		case "email":
			email, err = d.Str()
		case "username":
			username, err = d.Str()
		case "password":
			password, err = d.Str()
		case "provider":
			provider, err = d.Str()
		case "superuser":
			superuser, err = d.Str()
		default:
			return errors.New("Bad Request")
		}
		return err
	})
	
	if err != nil {
		logger.Error().Err(err).Msg("User creation failed")
		BadRequest.Write(res)
		return
	}

	if fname == "" || lname == "" || email == "" || username == "" || provider == "" {
		logger.Debug().
			Str("fname", fname).
			Str("lname", lname).
			Str("email", email).
			Str("username", username).
			Str("provider", provider).
			Msg("Request validation failed.")
		BadRequest.Write(res)
		return
	}

	var providerType usr.ProviderType
	switch provider {
	case "LDAP", "ldap":
		providerType = usr.LDAP
	case "LOCAL", "local":
		providerType = usr.LOCAL
	default:
		logger.Error().Str("provider", provider).Msg("Unsupported Provider Type")
		BadRequest.Write(res)
		return
	}

	if err := h.Context.UserProvider.AddUser(providerType, fname, lname, email, username, password, superuser == "yes") ; err != nil {
		BadRequest.WriteWithError(res, err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("{\"message\":\"User Created\"}"))
}

func (h UserHandler) handleUpdate(res http.ResponseWriter, req *http.Request) {
	var fname, lname, email, username, password, provider string

	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "fname":
			fname, err = d.Str()
		case "lname":
			lname, err = d.Str()
		case "email":
			email, err = d.Str()
		case "username":
			username, err = d.Str()
		case "password":
			password, err = d.Str()
		case "provider":
			provider, err = d.Str()
		default:
			return errors.New("Bad Request")
		}
		return err
	})
	
	if err != nil {
		logger.Error().Err(err).Msg("User creation failed")
		BadRequest.Write(res)
		return
	}

	if fname == "" && lname == "" && email == "" && password == "" {
		logger.Debug().
			Str("fname", fname).
			Str("lname", lname).
			Str("email", email).
			Str("username", username).
			Str("provider", provider).
			Msg("Request validation failed.")
		BadRequest.Write(res)
		return
	}

	var providerType usr.ProviderType
	switch provider {
	case "LDAP", "ldap":
		providerType = usr.LDAP
	case "LOCAL", "local":
		providerType = usr.LOCAL
	default:
		logger.Error().Str("provider", provider).Msg("Unsupported Provider Type")
		BadRequest.Write(res)
		return
	}

	if err := h.Context.UserProvider.UpdateUser(providerType, fname, lname, email, username, password); err != nil {
		logger.Error().Err(err).Msg("Failed to update user")
		BadRequest.WriteWithError(res, err)
		return
	}

	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte("{\"message\":\"User Updated\"}"))
}
