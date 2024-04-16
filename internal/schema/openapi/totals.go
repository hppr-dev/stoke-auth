package openapi

import "github.com/ogen-go/ogen"

func addTotalsEndpoint(spec *ogen.Spec, security ogen.SecurityRequirements) error {
	pathItem := ogen.NewPathItem().
		SetDescription("Entity count totals").
		SetGet(ogen.NewOperation().
			SetOperationID("totals").
			SetSummary("Get entity count totals").
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().SetName("users").SetSchema(ogen.Int().SetDescription("Count of users")),
						*ogen.NewProperty().SetName("claims").SetSchema(ogen.Int().SetDescription("Count of claims")),
						*ogen.NewProperty().SetName("claim_groups").SetSchema(ogen.Int().SetDescription("Count of claim groups")),
					}).
					SetRequired([]string{"users", "claims", "claim_groups"}),
				),
			),
		)
	pathItem.Get.Security = security
	spec.AddPathItem("/admin/totals", pathItem)
	return nil
}
