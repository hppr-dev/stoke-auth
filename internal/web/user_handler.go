package web

import (
	"errors"
	"net/http"
	"stoke/internal/tel"
	"stoke/internal/usr"

	"github.com/go-faster/jx"
	"github.com/rs/zerolog"
)

type UserHandler struct {}

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
	ctx := req.Context()
	logger := zerolog.Ctx(ctx)
	_, span := tel.GetTracer().Start(ctx, "UserHandler.handleCreate")
	defer span.End()

	var fname, lname, email, username, password, superuser string

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

	if fname == "" || lname == "" || email == "" || username == "" {
		logger.Debug().
			Str("fname", fname).
			Str("lname", lname).
			Str("email", email).
			Str("username", username).
			Msg("Request validation failed.")
		BadRequest.Write(res)
		return
	}

	if err := usr.ProviderFromCtx(ctx).AddUser(fname, lname, email, username, password, superuser == "yes", ctx) ; err != nil {
		BadRequest.WriteWithError(res, err)
		return
	}
	res.WriteHeader(http.StatusCreated)
	res.Write([]byte("{\"message\":\"User Created\"}"))
}

func (h UserHandler) handleUpdate(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	logger := zerolog.Ctx(ctx)
	_, span := tel.GetTracer().Start(ctx, "UserHandler.handleUpdate")
	defer span.End()

	var username, newPassword, oldPassword string
	var force bool

	decoder := jx.Decode(req.Body, 256)
	err := decoder.Obj(func (d *jx.Decoder, key string) error {
		var err error
		switch key {
		case "username":
			username, err = d.Str()
		case "oldPassword":
			oldPassword, err = d.Str()
		case "newPassword":
			newPassword, err = d.Str()
		case "force":
			force, err = d.Bool()
		default:
			return errors.New("Bad Request")
		}
		return err
	})
	
	if err != nil {
		logger.Error().Err(err).Msg("User password update failed")
		BadRequest.Write(res)
		return
	}

	if newPassword == "" {
		logger.Debug().
			Msg("Request validation failed.")
		BadRequest.Write(res)
		return
	}

	if err := usr.ProviderFromCtx(ctx).UpdateUserPassword(username, oldPassword, newPassword, force, ctx); err != nil {
		logger.Error().Err(err).Msg("Failed to update user")
		BadRequest.WriteWithError(res, err)
		return
	}

	res.WriteHeader(http.StatusAccepted)
	res.Write([]byte("{\"message\":\"User Updated\"}"))
}
