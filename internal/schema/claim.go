package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Claim struct {
	ent.Schema
}

func (Claim) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("short_name").
			Unique(),
		field.String("description"),
	}
}

func (Claim) Edges() []ent.Edge {
	return []ent.Edge {
		edge.To("claim_groups", ClaimGroup.Type),
	}
}

func (Claim) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
