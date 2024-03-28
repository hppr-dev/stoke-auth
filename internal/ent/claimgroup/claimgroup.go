// Code generated by ent, DO NOT EDIT.

package claimgroup

import (
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the claimgroup type in the database.
	Label = "claim_group"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldDescription holds the string denoting the description field in the database.
	FieldDescription = "description"
	// FieldIsUserGroup holds the string denoting the is_user_group field in the database.
	FieldIsUserGroup = "is_user_group"
	// EdgeUsers holds the string denoting the users edge name in mutations.
	EdgeUsers = "users"
	// EdgeGroupLinks holds the string denoting the group_links edge name in mutations.
	EdgeGroupLinks = "group_links"
	// EdgeClaims holds the string denoting the claims edge name in mutations.
	EdgeClaims = "claims"
	// Table holds the table name of the claimgroup in the database.
	Table = "claim_groups"
	// UsersTable is the table that holds the users relation/edge. The primary key declared below.
	UsersTable = "claim_group_users"
	// UsersInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UsersInverseTable = "users"
	// GroupLinksTable is the table that holds the group_links relation/edge.
	GroupLinksTable = "group_links"
	// GroupLinksInverseTable is the table name for the GroupLink entity.
	// It exists in this package in order to avoid circular dependency with the "grouplink" package.
	GroupLinksInverseTable = "group_links"
	// GroupLinksColumn is the table column denoting the group_links relation/edge.
	GroupLinksColumn = "claim_group_group_links"
	// ClaimsTable is the table that holds the claims relation/edge. The primary key declared below.
	ClaimsTable = "claim_claim_groups"
	// ClaimsInverseTable is the table name for the Claim entity.
	// It exists in this package in order to avoid circular dependency with the "claim" package.
	ClaimsInverseTable = "claims"
)

// Columns holds all SQL columns for claimgroup fields.
var Columns = []string{
	FieldID,
	FieldName,
	FieldDescription,
	FieldIsUserGroup,
}

var (
	// UsersPrimaryKey and UsersColumn2 are the table columns denoting the
	// primary key for the users relation (M2M).
	UsersPrimaryKey = []string{"claim_group_id", "user_id"}
	// ClaimsPrimaryKey and ClaimsColumn2 are the table columns denoting the
	// primary key for the claims relation (M2M).
	ClaimsPrimaryKey = []string{"claim_id", "claim_group_id"}
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

var (
	// DefaultIsUserGroup holds the default value on creation for the "is_user_group" field.
	DefaultIsUserGroup bool
)

// OrderOption defines the ordering options for the ClaimGroup queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByDescription orders the results by the description field.
func ByDescription(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDescription, opts...).ToFunc()
}

// ByIsUserGroup orders the results by the is_user_group field.
func ByIsUserGroup(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsUserGroup, opts...).ToFunc()
}

// ByUsersCount orders the results by users count.
func ByUsersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newUsersStep(), opts...)
	}
}

// ByUsers orders the results by users terms.
func ByUsers(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUsersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByGroupLinksCount orders the results by group_links count.
func ByGroupLinksCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newGroupLinksStep(), opts...)
	}
}

// ByGroupLinks orders the results by group_links terms.
func ByGroupLinks(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newGroupLinksStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByClaimsCount orders the results by claims count.
func ByClaimsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newClaimsStep(), opts...)
	}
}

// ByClaims orders the results by claims terms.
func ByClaims(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newClaimsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newUsersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UsersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, UsersTable, UsersPrimaryKey...),
	)
}
func newGroupLinksStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(GroupLinksInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, GroupLinksTable, GroupLinksColumn),
	)
}
func newClaimsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ClaimsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, ClaimsTable, ClaimsPrimaryKey...),
	)
}
