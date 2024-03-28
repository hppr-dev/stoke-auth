// Code generated by ent, DO NOT EDIT.

package claimgroup

import (
	"stoke/internal/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLTE(FieldID, id))
}

// Name applies equality check predicate on the "name" field. It's identical to NameEQ.
func Name(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldName, v))
}

// Description applies equality check predicate on the "description" field. It's identical to DescriptionEQ.
func Description(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldDescription, v))
}

// IsUserGroup applies equality check predicate on the "is_user_group" field. It's identical to IsUserGroupEQ.
func IsUserGroup(v bool) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldIsUserGroup, v))
}

// NameEQ applies the EQ predicate on the "name" field.
func NameEQ(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldName, v))
}

// NameNEQ applies the NEQ predicate on the "name" field.
func NameNEQ(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNEQ(FieldName, v))
}

// NameIn applies the In predicate on the "name" field.
func NameIn(vs ...string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldIn(FieldName, vs...))
}

// NameNotIn applies the NotIn predicate on the "name" field.
func NameNotIn(vs ...string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNotIn(FieldName, vs...))
}

// NameGT applies the GT predicate on the "name" field.
func NameGT(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGT(FieldName, v))
}

// NameGTE applies the GTE predicate on the "name" field.
func NameGTE(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGTE(FieldName, v))
}

// NameLT applies the LT predicate on the "name" field.
func NameLT(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLT(FieldName, v))
}

// NameLTE applies the LTE predicate on the "name" field.
func NameLTE(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLTE(FieldName, v))
}

// NameContains applies the Contains predicate on the "name" field.
func NameContains(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldContains(FieldName, v))
}

// NameHasPrefix applies the HasPrefix predicate on the "name" field.
func NameHasPrefix(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldHasPrefix(FieldName, v))
}

// NameHasSuffix applies the HasSuffix predicate on the "name" field.
func NameHasSuffix(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldHasSuffix(FieldName, v))
}

// NameEqualFold applies the EqualFold predicate on the "name" field.
func NameEqualFold(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEqualFold(FieldName, v))
}

// NameContainsFold applies the ContainsFold predicate on the "name" field.
func NameContainsFold(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldContainsFold(FieldName, v))
}

// DescriptionEQ applies the EQ predicate on the "description" field.
func DescriptionEQ(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldDescription, v))
}

// DescriptionNEQ applies the NEQ predicate on the "description" field.
func DescriptionNEQ(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNEQ(FieldDescription, v))
}

// DescriptionIn applies the In predicate on the "description" field.
func DescriptionIn(vs ...string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldIn(FieldDescription, vs...))
}

// DescriptionNotIn applies the NotIn predicate on the "description" field.
func DescriptionNotIn(vs ...string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNotIn(FieldDescription, vs...))
}

// DescriptionGT applies the GT predicate on the "description" field.
func DescriptionGT(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGT(FieldDescription, v))
}

// DescriptionGTE applies the GTE predicate on the "description" field.
func DescriptionGTE(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldGTE(FieldDescription, v))
}

// DescriptionLT applies the LT predicate on the "description" field.
func DescriptionLT(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLT(FieldDescription, v))
}

// DescriptionLTE applies the LTE predicate on the "description" field.
func DescriptionLTE(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldLTE(FieldDescription, v))
}

// DescriptionContains applies the Contains predicate on the "description" field.
func DescriptionContains(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldContains(FieldDescription, v))
}

// DescriptionHasPrefix applies the HasPrefix predicate on the "description" field.
func DescriptionHasPrefix(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldHasPrefix(FieldDescription, v))
}

// DescriptionHasSuffix applies the HasSuffix predicate on the "description" field.
func DescriptionHasSuffix(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldHasSuffix(FieldDescription, v))
}

// DescriptionEqualFold applies the EqualFold predicate on the "description" field.
func DescriptionEqualFold(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEqualFold(FieldDescription, v))
}

// DescriptionContainsFold applies the ContainsFold predicate on the "description" field.
func DescriptionContainsFold(v string) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldContainsFold(FieldDescription, v))
}

// IsUserGroupEQ applies the EQ predicate on the "is_user_group" field.
func IsUserGroupEQ(v bool) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldEQ(FieldIsUserGroup, v))
}

// IsUserGroupNEQ applies the NEQ predicate on the "is_user_group" field.
func IsUserGroupNEQ(v bool) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.FieldNEQ(FieldIsUserGroup, v))
}

// HasUsers applies the HasEdge predicate on the "users" edge.
func HasUsers() predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, false, UsersTable, UsersPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUsersWith applies the HasEdge predicate on the "users" edge with a given conditions (other predicates).
func HasUsersWith(preds ...predicate.User) predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := newUsersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasGroupLinks applies the HasEdge predicate on the "group_links" edge.
func HasGroupLinks() predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, GroupLinksTable, GroupLinksColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasGroupLinksWith applies the HasEdge predicate on the "group_links" edge with a given conditions (other predicates).
func HasGroupLinksWith(preds ...predicate.GroupLink) predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := newGroupLinksStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasClaims applies the HasEdge predicate on the "claims" edge.
func HasClaims() predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ClaimsTable, ClaimsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasClaimsWith applies the HasEdge predicate on the "claims" edge with a given conditions (other predicates).
func HasClaimsWith(preds ...predicate.Claim) predicate.ClaimGroup {
	return predicate.ClaimGroup(func(s *sql.Selector) {
		step := newClaimsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ClaimGroup) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ClaimGroup) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ClaimGroup) predicate.ClaimGroup {
	return predicate.ClaimGroup(sql.NotPredicates(p))
}
