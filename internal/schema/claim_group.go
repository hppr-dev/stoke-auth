package schema

import (
	"stoke/internal/ent/privacy"
	"stoke/internal/ent/schema/policy"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
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
		edge.To("group_links", GroupLink.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
		edge.From("claims", Claim.Type).
			Ref("claim_groups"),
	}
}

func (ClaimGroup) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}

func (ClaimGroup) Policy() ent.Policy {
	return privacy.Policy {
		Mutation: privacy.MutationPolicy{
			policy.ClaimGroupMutationPolicy{},
		},
	}
}
