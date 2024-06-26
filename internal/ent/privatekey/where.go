// Code generated by ent, DO NOT EDIT.

package privatekey

import (
	"stoke/internal/ent/predicate"
	"time"

	"entgo.io/ent/dialect/sql"
)

// ID filters vertices based on their ID field.
func ID(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id int) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLTE(FieldID, id))
}

// Text applies equality check predicate on the "text" field. It's identical to TextEQ.
func Text(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldText, v))
}

// Expires applies equality check predicate on the "expires" field. It's identical to ExpiresEQ.
func Expires(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldExpires, v))
}

// TextEQ applies the EQ predicate on the "text" field.
func TextEQ(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldText, v))
}

// TextNEQ applies the NEQ predicate on the "text" field.
func TextNEQ(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNEQ(FieldText, v))
}

// TextIn applies the In predicate on the "text" field.
func TextIn(vs ...string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldIn(FieldText, vs...))
}

// TextNotIn applies the NotIn predicate on the "text" field.
func TextNotIn(vs ...string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNotIn(FieldText, vs...))
}

// TextGT applies the GT predicate on the "text" field.
func TextGT(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGT(FieldText, v))
}

// TextGTE applies the GTE predicate on the "text" field.
func TextGTE(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGTE(FieldText, v))
}

// TextLT applies the LT predicate on the "text" field.
func TextLT(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLT(FieldText, v))
}

// TextLTE applies the LTE predicate on the "text" field.
func TextLTE(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLTE(FieldText, v))
}

// TextContains applies the Contains predicate on the "text" field.
func TextContains(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldContains(FieldText, v))
}

// TextHasPrefix applies the HasPrefix predicate on the "text" field.
func TextHasPrefix(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldHasPrefix(FieldText, v))
}

// TextHasSuffix applies the HasSuffix predicate on the "text" field.
func TextHasSuffix(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldHasSuffix(FieldText, v))
}

// TextEqualFold applies the EqualFold predicate on the "text" field.
func TextEqualFold(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEqualFold(FieldText, v))
}

// TextContainsFold applies the ContainsFold predicate on the "text" field.
func TextContainsFold(v string) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldContainsFold(FieldText, v))
}

// ExpiresEQ applies the EQ predicate on the "expires" field.
func ExpiresEQ(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldEQ(FieldExpires, v))
}

// ExpiresNEQ applies the NEQ predicate on the "expires" field.
func ExpiresNEQ(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNEQ(FieldExpires, v))
}

// ExpiresIn applies the In predicate on the "expires" field.
func ExpiresIn(vs ...time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldIn(FieldExpires, vs...))
}

// ExpiresNotIn applies the NotIn predicate on the "expires" field.
func ExpiresNotIn(vs ...time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldNotIn(FieldExpires, vs...))
}

// ExpiresGT applies the GT predicate on the "expires" field.
func ExpiresGT(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGT(FieldExpires, v))
}

// ExpiresGTE applies the GTE predicate on the "expires" field.
func ExpiresGTE(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldGTE(FieldExpires, v))
}

// ExpiresLT applies the LT predicate on the "expires" field.
func ExpiresLT(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLT(FieldExpires, v))
}

// ExpiresLTE applies the LTE predicate on the "expires" field.
func ExpiresLTE(v time.Time) predicate.PrivateKey {
	return predicate.PrivateKey(sql.FieldLTE(FieldExpires, v))
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.PrivateKey) predicate.PrivateKey {
	return predicate.PrivateKey(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.PrivateKey) predicate.PrivateKey {
	return predicate.PrivateKey(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.PrivateKey) predicate.PrivateKey {
	return predicate.PrivateKey(sql.NotPredicates(p))
}
