package schema

import (
	"time"

	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
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
					Annotations(entoas.ListOperation(entoas.OperationPolicy(entoas.PolicyExclude))),
				field.String("salt").
					Annotations(entoas.ListOperation(entoas.OperationPolicy(entoas.PolicyExclude))),
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
		entoas.DeleteOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
	}
}

func (User) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
