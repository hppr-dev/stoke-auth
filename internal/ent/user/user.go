// Code generated by ent, DO NOT EDIT.

package user

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the user type in the database.
	Label = "user"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldFname holds the string denoting the fname field in the database.
	FieldFname = "fname"
	// FieldLname holds the string denoting the lname field in the database.
	FieldLname = "lname"
	// FieldSource holds the string denoting the source field in the database.
	FieldSource = "source"
	// FieldEmail holds the string denoting the email field in the database.
	FieldEmail = "email"
	// FieldUsername holds the string denoting the username field in the database.
	FieldUsername = "username"
	// FieldPassword holds the string denoting the password field in the database.
	FieldPassword = "password"
	// FieldSalt holds the string denoting the salt field in the database.
	FieldSalt = "salt"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeClaimGroups holds the string denoting the claim_groups edge name in mutations.
	EdgeClaimGroups = "claim_groups"
	// Table holds the table name of the user in the database.
	Table = "users"
	// ClaimGroupsTable is the table that holds the claim_groups relation/edge. The primary key declared below.
	ClaimGroupsTable = "claim_group_users"
	// ClaimGroupsInverseTable is the table name for the ClaimGroup entity.
	// It exists in this package in order to avoid circular dependency with the "claimgroup" package.
	ClaimGroupsInverseTable = "claim_groups"
)

// Columns holds all SQL columns for user fields.
var Columns = []string{
	FieldID,
	FieldFname,
	FieldLname,
	FieldSource,
	FieldEmail,
	FieldUsername,
	FieldPassword,
	FieldSalt,
	FieldCreatedAt,
}

var (
	// ClaimGroupsPrimaryKey and ClaimGroupsColumn2 are the table columns denoting the
	// primary key for the claim_groups relation (M2M).
	ClaimGroupsPrimaryKey = []string{"claim_group_id", "user_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	return false
}

// Note that the variables below are initialized by the runtime
// package on the initialization of the application. Therefore,
// it should be imported in the main as follows:
//
//	import _ "stoke/internal/ent/runtime"
var (
	Hooks  [1]ent.Hook
	Policy ent.Policy
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
)

// OrderOption defines the ordering options for the User queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByFname orders the results by the fname field.
func ByFname(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFname, opts...).ToFunc()
}

// ByLname orders the results by the lname field.
func ByLname(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldLname, opts...).ToFunc()
}

// BySource orders the results by the source field.
func BySource(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSource, opts...).ToFunc()
}

// ByEmail orders the results by the email field.
func ByEmail(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldEmail, opts...).ToFunc()
}

// ByUsername orders the results by the username field.
func ByUsername(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUsername, opts...).ToFunc()
}

// ByPassword orders the results by the password field.
func ByPassword(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldPassword, opts...).ToFunc()
}

// BySalt orders the results by the salt field.
func BySalt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSalt, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByClaimGroupsCount orders the results by claim_groups count.
func ByClaimGroupsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newClaimGroupsStep(), opts...)
	}
}

// ByClaimGroups orders the results by claim_groups terms.
func ByClaimGroups(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newClaimGroupsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newClaimGroupsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ClaimGroupsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, ClaimGroupsTable, ClaimGroupsPrimaryKey...),
	)
}
