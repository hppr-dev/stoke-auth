package openapi

import "github.com/ogen-go/ogen"

func addCapabilitesEndpoint(spec *ogen.Spec, security ogen.SecurityRequirements) error {
	pathItem := ogen.NewPathItem().
		SetDescription("Server capabilities").
		SetGet(ogen.NewOperation().
			SetOperationID("capabilities").
			SetSummary("Get server capabilities").
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("capabilities").SetSchema(ogen.String().AsArray().SetDescription("List of enabled capabilites")),
					}).
					SetRequired([]string{"capabilities"}),
				),
			),
		)
	pathItem.Get.Security = security
	spec.AddPathItem("/capabilities", pathItem)
	return nil
}
