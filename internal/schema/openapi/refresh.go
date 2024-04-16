package openapi

import (
	"encoding/json"

	"github.com/ogen-go/ogen"
)

func addRefreshEndpoint(spec *ogen.Spec, security ogen.SecurityRequirements) error {
	pathItem := ogen.NewPathItem().
		SetDescription("Token refresh endpoint").
		SetPost(ogen.NewOperation().
			SetOperationID("refresh").
			SetSummary("Request a refreshed token").
			SetRequestBody(ogen.NewRequestBody().
				SetRequired(true).
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetDescription("User credentials").
					SetRequired([]string{"refresh"}).
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("refresh").
							SetSchema(ogen.String().
								SetDescription("Refresh token. Must match the token used in authentication"),
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
					}).
					SetRequired([]string{"token", "refresh"}),
				),
			),
		)
	pathItem.Post.Security = security
	spec.AddPathItem("/refresh", pathItem)
	return nil
}
