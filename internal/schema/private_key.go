package schema

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type PrivateKey struct {
		ent.Schema
}

func (PrivateKey) Fields() []ent.Field {
		return []ent.Field{
			field.String("text").
				Immutable(),
			field.Time("expires").
				Immutable(),
			field.Time("renews").
				Immutable(),
		}
}

func (PrivateKey) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entoas.CreateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.DeleteOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.UpdateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
	}
}

func (PrivateKey) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}
