// Code generated by ent, DO NOT EDIT.

package privatekey

import (
	"entgo.io/ent/dialect/sql"
)

const (
	// Label holds the string label denoting the privatekey type in the database.
	Label = "private_key"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldText holds the string denoting the text field in the database.
	FieldText = "text"
	// FieldExpires holds the string denoting the expires field in the database.
	FieldExpires = "expires"
	// FieldRenews holds the string denoting the renews field in the database.
	FieldRenews = "renews"
	// Table holds the table name of the privatekey in the database.
	Table = "private_keys"
)

// Columns holds all SQL columns for privatekey fields.
var Columns = []string{
	FieldID,
	FieldText,
	FieldExpires,
	FieldRenews,
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// OrderOption defines the ordering options for the PrivateKey queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByText orders the results by the text field.
func ByText(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldText, opts...).ToFunc()
}

// ByExpires orders the results by the expires field.
func ByExpires(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldExpires, opts...).ToFunc()
}

// ByRenews orders the results by the renews field.
func ByRenews(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRenews, opts...).ToFunc()
}
