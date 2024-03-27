package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type ClaimGroup struct {
	ent.Schema
}

func (ClaimGroup) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("description"),
	}
}

func (ClaimGroup) Edges() []ent.Edge {
	return []ent.Edge {
		edge.To("users", User.Type),
		edge.To("group_links", GroupLink.Type),
		edge.From("claims", Claim.Type).
			Ref("claim_groups"),
	}
}

func (ClaimGroup) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
