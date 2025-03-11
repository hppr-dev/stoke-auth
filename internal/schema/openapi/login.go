package openapi

import (
	"encoding/json"

	"github.com/ogen-go/ogen"
)

func addLoginEndpoint(spec *ogen.Spec) error {
	mapObjectSchema := ogen.NewSchema().SetType("object")
	mapObjectSchema.AdditionalProperties = &ogen.AdditionalProperties{ Schema : *ogen.String() }

	pathItem := ogen.NewPathItem().
		SetDescription("User login and token generation endpoint").
		SetPost(ogen.NewOperation().
			SetOperationID("login").
			SetSummary("Request a token").
			SetRequestBody(ogen.NewRequestBody().
				SetRequired(true).
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetDescription("User credentials").
					SetRequired([]string{"username", "password"}).
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("username").
							SetSchema(ogen.String().
								SetDescription("User's username"),
							),
						*ogen.NewProperty().
							SetName("password").
							SetSchema(ogen.String().
								SetDescription("User's password"),
							),
						*ogen.NewProperty().
							SetName("provider").
							SetSchema(ogen.String().
								SetDescription("Provider to login try to login to. This is required when multiple foreign providers are defined."),
							),
						*ogen.NewProperty().
							SetName("required_claims").
							SetSchema(ogen.NewSchema().
								SetType("array").
								SetDescription("Claims required to receive a token. A token is issued if a user matches all given claims of one entry in the list.").
								SetItems(mapObjectSchema.SetDescription("An object that specifies what claims must be present. Set a key value to \"\" to require a key with any value"),
								),
							),
						*ogen.NewProperty().
							SetName("filter_claims").
							SetSchema(ogen.NewSchema().
								SetType("array").
								SetDescription("Claims to include in the created token. All given claims are returned by default").
								SetItems(ogen.String()),
							),
					}),
				),
			).
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("token").
							SetSchema(ogen.String().
								SetDescription("JWT Token"),
							),
						*ogen.NewProperty().
							SetName("refresh").
							SetSchema(ogen.String().
								SetDescription("Token to get a new token with the same claims. Must be used before token expires"),
							),
						*ogen.NewProperty().
							SetName("username").
							SetSchema(ogen.String().
								SetDescription("Username of the user who logged in"),
							),
					}).
					SetRequired([]string{"username", "token", "refresh"}),
				),
			).
			AddResponse("401", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("message").
							SetSchema(ogen.String().
								SetDescription("Error Message").
								SetDefault(json.RawMessage(`"Not Authorized"`)),
							),
					}),
				),
			).
			AddResponse("400", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("message").
							SetSchema(ogen.String().
								SetDescription("Error Message").
								SetDefault(json.RawMessage(`"Unprocessable Entry"`)),
							),
					}),
				),
			),
		)
	spec.AddPathItem("/login", pathItem)
	return nil
}
