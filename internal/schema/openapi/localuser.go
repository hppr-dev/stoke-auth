package openapi

import (
	"encoding/json"

	"github.com/ogen-go/ogen"
)

func addLocalUserEndpoint(spec *ogen.Spec, security ogen.SecurityRequirements) error {
	pathItem := ogen.NewPathItem().
		SetDescription("Local user modification functions").
		SetPost(ogen.NewOperation().
			SetOperationID("createLocalUser").
			SetSummary("Create a new local user").
			SetRequestBody(ogen.NewRequestBody().
				SetRequired(true).
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetDescription("User to create").
					SetRequired([]string{"fname", "lname", "username", "email", "password"}).
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("fname").SetSchema(ogen.String().SetDescription("New user's first name")),
						*ogen.NewProperty().SetName("lname").SetSchema(ogen.String().SetDescription("New user's last name")),
						*ogen.NewProperty().SetName("username").SetSchema(ogen.String().SetDescription("New user's username")),
						*ogen.NewProperty().SetName("email").SetSchema(ogen.String().SetDescription("New user's email")),
						*ogen.NewProperty().SetName("password").SetSchema(ogen.String().SetDescription("New user's password")),
					}),
				),
			).
			AddResponse("400", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("message").SetSchema(ogen.String().SetDescription("Error Message")),
					}).
					SetRequired([]string{"message"}),
				),
			).
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("message").SetSchema(ogen.String().SetDefault(json.RawMessage(`"User Created"`))),
					}),
				),
			),
		).
		SetPatch(ogen.NewOperation().
			SetOperationID("updateLocalUserPassword").
			SetSummary("Update local user's password").
			SetRequestBody(ogen.NewRequestBody().
				SetRequired(true).
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetDescription("Password update data").
					SetRequired([]string{"username", "newPassword"}).
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("username").SetSchema(ogen.String().SetDescription("User's username")),
						*ogen.NewProperty().SetName("oldPassword").SetSchema(ogen.String().SetDescription("Old user password")),
						*ogen.NewProperty().SetName("newPassword").SetSchema(ogen.String().SetDescription("New user password")),
						*ogen.NewProperty().SetName("force").SetSchema(ogen.Bool().SetDescription("Set to change user's password without checking the old password")),
					}),
				),
			).
			AddResponse("400", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("message").SetSchema(ogen.String().SetDescription("Error Message")),
					}).
					SetRequired([]string{"message"}),
				),
			).
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("message").SetSchema(ogen.String().SetDefault(json.RawMessage(`"Password Updated"`))),
					}),
				),
			),
		)
	pathItem.Post.Security = security
	pathItem.Patch.Security = security
	spec.AddPathItem("/admin/localuser", pathItem)
	return nil
}
