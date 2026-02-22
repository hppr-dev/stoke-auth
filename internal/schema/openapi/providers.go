package openapi

import "github.com/ogen-go/ogen"

func addAvailableProvidersEndpoint(spec *ogen.Spec) error {
	pathItem := ogen.NewPathItem().
		SetDescription("Lists which providers are available").
		SetGet(ogen.NewOperation().
			SetOperationID("available_providers").
			SetSummary("Get available providers").
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("providers").SetSchema(ogen.NewSchema().
							SetType("array").
							SetItems(ogen.NewSchema().
								SetType("object").
								SetProperties(&ogen.Properties{
									*ogen.NewProperty().SetName("name").SetSchema(ogen.String().SetDescription("Name of provider")),
									*ogen.NewProperty().SetName("provider_type").SetSchema(ogen.String().SetDescription("Type of provider")),
									*ogen.NewProperty().SetName("type_spec").SetSchema(ogen.String().SetDescription("Type specification of provider")),
								}).
								SetRequired([]string{"name", "provider_type", "type_spec"}),
							),
						),
						*ogen.NewProperty().SetName("base_admin_path").SetSchema(ogen.String().SetDescription("Base path for the admin UI when served behind a proxy, e.g. /auth; empty when admin is at /admin/")),
					}).
					SetRequired([]string{"providers"}),
				),
			),
		)
	spec.AddPathItem("/available_providers", pathItem)
	return nil
}
