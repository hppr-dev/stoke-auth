package schema

import (
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

type DBInitFile struct {
	ent.Schema
}

func (DBInitFile) Fields() []ent.Field {
	return []ent.Field{
		field.String("filename"),
		field.String("md5"),
	}
}

func (DBInitFile) Mixins() []ent.Mixin {
	return []ent.Mixin{
		Common{},
	}
}

func (DBInitFile) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entoas.CreateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.ReadOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.UpdateOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.DeleteOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
		entoas.ListOperation(entoas.OperationPolicy(entoas.PolicyExclude)),
	}
}
