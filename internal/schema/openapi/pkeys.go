package openapi

import (
	"encoding/json"

	"github.com/ogen-go/ogen"
)

func addPkeysEndpoint(spec *ogen.Spec) error {
	spec.AddPathItem("/pkeys", ogen.NewPathItem().
		SetDescription("Current Public keys").
		SetGet(ogen.NewOperation().
			SetOperationID("pkeys").
			SetSummary("Get current valid public keys").
			SetDescription("Returns JWKS (merged from all peers when clustered). Optional query: local=true or local=1 to return only this node's keys (used by peers to avoid recursion).").
			AddResponse("200", ogen.NewResponse().
				AddContent("application/json", ogen.NewSchema().
					SetType("object").
					SetProperties(&ogen.Properties{
						*ogen.NewProperty().
							SetName("exp").
							SetSchema(ogen.DateTime().
								SetDescription("Next key expiring time, next time to pull public keys"),
							),
						*ogen.NewProperty().
							SetName("keys").
							SetSchema(ogen.NewSchema().
								SetType("array").
								SetItems(ogen.NewSchema().
									SetType("object").
									SetProperties(&ogen.Properties{
										*ogen.NewProperty().SetName("kty").SetSchema(
											ogen.String().SetEnum([]json.RawMessage{
												json.RawMessage(`"EC"`),
												json.RawMessage(`"RSA"`),
												json.RawMessage(`"OKP"`),
											}).SetDescription("Key Type")),
										*ogen.NewProperty().SetName("use").SetSchema(ogen.String().SetDefault(json.RawMessage(`"sig"`)).SetDescription("Key usage")),
										*ogen.NewProperty().SetName("kid").SetSchema(ogen.String().SetDescription("Key identifier")),

										*ogen.NewProperty().SetName("crv").SetSchema(
											ogen.String().SetEnum([]json.RawMessage{
												json.RawMessage(`"P-256"`),
												json.RawMessage(`"P-384"`),
												json.RawMessage(`"P-521"`),
											}).SetDescription("ECDSA/EdDSA Curve")),
										*ogen.NewProperty().SetName("x").SetSchema(ogen.String().SetDescription("URL encoded base64 ECDSA/EdDSA X")),
										*ogen.NewProperty().SetName("y").SetSchema(ogen.String().SetDescription("URL encoded base64 ECDSA Y")),

										*ogen.NewProperty().SetName("n").SetSchema(ogen.String().SetDescription("URL encoded base64 RSA N")),
										*ogen.NewProperty().SetName("e").SetSchema(ogen.String().SetDescription("URL encoded base64 RSA E")),
									}),
								),
							),
						}).
						SetRequired([]string{"exp", "keys"}),
					),
				),
			),
	)
	return nil
}
