package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type GroupLink struct {
	ent.Schema
}

func (GroupLink) Fields() []ent.Field {
	return []ent.Field{
		field.String("type"),
		field.String("resource_spec"),
	}
}

func (GroupLink) Edges() []ent.Edge {
	return []ent.Edge {
		edge.From("claim_groups", ClaimGroup.Type).
			Ref("group_links").
			Unique(),
	}
}

func (GroupLink) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
