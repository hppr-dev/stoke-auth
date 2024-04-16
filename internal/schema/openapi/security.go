package openapi

import (
	"github.com/ogen-go/ogen"
)

func addSecurity(spec *ogen.Spec) ogen.SecurityRequirements {
	spec.Components.SecuritySchemes = map[string]*ogen.SecurityScheme{
				"token" : {
					Type: "http",
					Description: "Requires Stoke Token Authentication",
					Scheme: "bearer",
				},
			}

	return ogen.SecurityRequirements{
			ogen.SecurityRequirement{
				"token" :  []string{},
			},
		}
}
