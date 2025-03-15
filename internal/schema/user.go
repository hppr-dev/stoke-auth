package schema

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"stoke/internal/ent/privacy"
	"stoke/internal/ent/schema/policy"
)

type User struct {
	ent.Schema
}

func (User) Fields() []ent.Field {
		return []ent.Field{
				field.String("fname"),
				field.String("lname"),
				field.String("source"),
				field.String("email").
					Unique(),
				field.String("username").
					Unique(),
				field.String("password").
					Optional().
					Annotations(
						entoas.Skip(true),
					),
				field.String("salt").
					Optional().
					Annotations(
						entoas.Skip(true),
					),
				field.Time("created_at").
					Immutable().
					Default(time.Now).
					Annotations(entoas.ReadOnly(true)),
		}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge {
		edge.From("claim_groups", ClaimGroup.Type).
			Ref("users"),
	}
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entoas.CreateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
	}
}

func (User) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}

func (User) Policy() ent.Policy {
	return privacy.Policy {
		Mutation: privacy.MutationPolicy{
			policy.UserMutationPolicy{},
		},
	}
}
