package schema

import (
	"time"

	"entgo.io/ent"
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
        field.String("email").
					Unique(),
        field.String("username").
        	Unique(),
        field.String("password"),
        field.String("salt"),
        field.Time("created_at").
            Default(time.Now),
    }
}

func (User) Edges() []ent.Edge {
	return []ent.Edge {
		edge.From("claim_groups", ClaimGroup.Type).
			Ref("users"),
	}
}

func (User) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
