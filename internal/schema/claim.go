package schema

import (
	"stoke/internal/ent/privacy"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

type Claim struct {
	ent.Schema
}

func (Claim) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique(),
		field.String("short_name"),
		field.String("value"),
		field.String("description"),
	}
}

func (Claim) Edges() []ent.Edge {
	return []ent.Edge {
		edge.To("claim_groups", ClaimGroup.Type),
	}
}

func (Claim) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("short_name", "value").
			Unique(),
	}
}

func (Claim) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}

func (Claim) Policy() ent.Policy {
	return privacy.Policy {
		Mutation: privacy.MutationPolicy{
			RestrictUpdates{
				EntityType: "claims",
				FieldName: "short_name",
			},
			privacy.AlwaysAllowRule(),
		},
	}
}
