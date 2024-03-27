package schema

import (
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/mixin"
)

type Common struct {
	mixin.Schema
}

func (Common) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Schema("stoke_auth"),
	}
}
