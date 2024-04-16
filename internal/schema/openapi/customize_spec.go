package openapi

import (
	"entgo.io/ent/entc/gen"
	"github.com/ogen-go/ogen"
)

func CustomizeSpec(g *gen.Graph, spec *ogen.Spec) error {
	security := addSecurity(spec)

	newPaths := make(ogen.Paths)
	for name, path := range spec.Paths {
		if path.Get != nil { path.Get.Security = security }
		if path.Post != nil { path.Post.Security = security }
		if path.Patch != nil { path.Patch.Security = security }
		if path.Delete != nil { path.Delete.Security = security }
		newPaths["/admin" + name] = path
	}
	spec.SetPaths(newPaths)

	addLocalUserEndpoint(spec, security)
	addTotalsEndpoint(spec, security)
	addRefreshEndpoint(spec, security)
	
	addLoginEndpoint(spec)
	addPkeysEndpoint(spec)
	return nil
}
