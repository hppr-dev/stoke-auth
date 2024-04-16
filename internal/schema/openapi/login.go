package openapi

import (
	"encoding/json"

	"github.com/ogen-go/ogen"
)

func addLoginEndpoint(spec *ogen.Spec) error {
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
							SetName("required_claims").
							SetSchema(ogen.NewSchema().
								SetType("array").
								SetDescription("Claims required to receive a token").
								SetItems(ogen.NewSchema().
									SetProperties(&ogen.Properties{
										*ogen.NewProperty().SetName("name").SetSchema(ogen.String().SetDescription("Required claim name")),
										*ogen.NewProperty().SetName("value").SetSchema(ogen.String().SetDescription("Required claim value")),
									}).
									SetRequired([]string{"name", "value"}),
								),
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
					}).
					SetRequired([]string{"token", "refresh"}),
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
