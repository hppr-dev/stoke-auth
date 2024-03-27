// Code generated by ent, DO NOT EDIT.

package grouplink

import (
	"stoke/internal/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLTE(FieldID, id))
}

// Type applies equality check predicate on the "type" field. It's identical to TypeEQ.
func Type(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldType, v))
}

// ResourceSpec applies equality check predicate on the "resource_spec" field. It's identical to ResourceSpecEQ.
func ResourceSpec(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldResourceSpec, v))
}

// TypeEQ applies the EQ predicate on the "type" field.
func TypeEQ(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldType, v))
}

// TypeNEQ applies the NEQ predicate on the "type" field.
func TypeNEQ(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNEQ(FieldType, v))
}

// TypeIn applies the In predicate on the "type" field.
func TypeIn(vs ...string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldIn(FieldType, vs...))
}

// TypeNotIn applies the NotIn predicate on the "type" field.
func TypeNotIn(vs ...string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNotIn(FieldType, vs...))
}

// TypeGT applies the GT predicate on the "type" field.
func TypeGT(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGT(FieldType, v))
}

// TypeGTE applies the GTE predicate on the "type" field.
func TypeGTE(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGTE(FieldType, v))
}

// TypeLT applies the LT predicate on the "type" field.
func TypeLT(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLT(FieldType, v))
}

// TypeLTE applies the LTE predicate on the "type" field.
func TypeLTE(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLTE(FieldType, v))
}

// TypeContains applies the Contains predicate on the "type" field.
func TypeContains(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldContains(FieldType, v))
}

// TypeHasPrefix applies the HasPrefix predicate on the "type" field.
func TypeHasPrefix(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldHasPrefix(FieldType, v))
}

// TypeHasSuffix applies the HasSuffix predicate on the "type" field.
func TypeHasSuffix(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldHasSuffix(FieldType, v))
}

// TypeEqualFold applies the EqualFold predicate on the "type" field.
func TypeEqualFold(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEqualFold(FieldType, v))
}

// TypeContainsFold applies the ContainsFold predicate on the "type" field.
func TypeContainsFold(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldContainsFold(FieldType, v))
}

// ResourceSpecEQ applies the EQ predicate on the "resource_spec" field.
func ResourceSpecEQ(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEQ(FieldResourceSpec, v))
}

// ResourceSpecNEQ applies the NEQ predicate on the "resource_spec" field.
func ResourceSpecNEQ(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNEQ(FieldResourceSpec, v))
}

// ResourceSpecIn applies the In predicate on the "resource_spec" field.
func ResourceSpecIn(vs ...string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldIn(FieldResourceSpec, vs...))
}

// ResourceSpecNotIn applies the NotIn predicate on the "resource_spec" field.
func ResourceSpecNotIn(vs ...string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldNotIn(FieldResourceSpec, vs...))
}

// ResourceSpecGT applies the GT predicate on the "resource_spec" field.
func ResourceSpecGT(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGT(FieldResourceSpec, v))
}

// ResourceSpecGTE applies the GTE predicate on the "resource_spec" field.
func ResourceSpecGTE(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldGTE(FieldResourceSpec, v))
}

// ResourceSpecLT applies the LT predicate on the "resource_spec" field.
func ResourceSpecLT(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLT(FieldResourceSpec, v))
}

// ResourceSpecLTE applies the LTE predicate on the "resource_spec" field.
func ResourceSpecLTE(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldLTE(FieldResourceSpec, v))
}

// ResourceSpecContains applies the Contains predicate on the "resource_spec" field.
func ResourceSpecContains(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldContains(FieldResourceSpec, v))
}

// ResourceSpecHasPrefix applies the HasPrefix predicate on the "resource_spec" field.
func ResourceSpecHasPrefix(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldHasPrefix(FieldResourceSpec, v))
}

// ResourceSpecHasSuffix applies the HasSuffix predicate on the "resource_spec" field.
func ResourceSpecHasSuffix(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldHasSuffix(FieldResourceSpec, v))
}

// ResourceSpecEqualFold applies the EqualFold predicate on the "resource_spec" field.
func ResourceSpecEqualFold(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldEqualFold(FieldResourceSpec, v))
}

// ResourceSpecContainsFold applies the ContainsFold predicate on the "resource_spec" field.
func ResourceSpecContainsFold(v string) predicate.GroupLink {
	return predicate.GroupLink(sql.FieldContainsFold(FieldResourceSpec, v))
}

// HasClaimGroups applies the HasEdge predicate on the "claim_groups" edge.
func HasClaimGroups() predicate.GroupLink {
	return predicate.GroupLink(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ClaimGroupsTable, ClaimGroupsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasClaimGroupsWith applies the HasEdge predicate on the "claim_groups" edge with a given conditions (other predicates).
func HasClaimGroupsWith(preds ...predicate.ClaimGroup) predicate.GroupLink {
	return predicate.GroupLink(func(s *sql.Selector) {
		step := newClaimGroupsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.GroupLink) predicate.GroupLink {
	return predicate.GroupLink(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.GroupLink) predicate.GroupLink {
	return predicate.GroupLink(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.GroupLink) predicate.GroupLink {
	return predicate.GroupLink(sql.NotPredicates(p))
}
