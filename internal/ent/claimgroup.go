// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"stoke/internal/ent/claimgroup"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// ClaimGroup is the model entity for the ClaimGroup schema.
type ClaimGroup struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ClaimGroupQuery when eager-loading is set.
	Edges        ClaimGroupEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ClaimGroupEdges holds the relations/edges for other nodes in the graph.
type ClaimGroupEdges struct {
	// Users holds the value of the users edge.
	Users []*User `json:"users,omitempty"`
	// GroupLinks holds the value of the group_links edge.
	GroupLinks []*GroupLink `json:"group_links,omitempty"`
	// Claims holds the value of the claims edge.
	Claims []*Claim `json:"claims,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// UsersOrErr returns the Users value or an error if the edge
// was not loaded in eager-loading.
func (e ClaimGroupEdges) UsersOrErr() ([]*User, error) {
	if e.loadedTypes[0] {
		return e.Users, nil
	}
	return nil, &NotLoadedError{edge: "users"}
}

// GroupLinksOrErr returns the GroupLinks value or an error if the edge
// was not loaded in eager-loading.
func (e ClaimGroupEdges) GroupLinksOrErr() ([]*GroupLink, error) {
	if e.loadedTypes[1] {
		return e.GroupLinks, nil
	}
	return nil, &NotLoadedError{edge: "group_links"}
}

// ClaimsOrErr returns the Claims value or an error if the edge
// was not loaded in eager-loading.
func (e ClaimGroupEdges) ClaimsOrErr() ([]*Claim, error) {
	if e.loadedTypes[2] {
		return e.Claims, nil
	}
	return nil, &NotLoadedError{edge: "claims"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ClaimGroup) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case claimgroup.FieldID:
			values[i] = new(sql.NullInt64)
		case claimgroup.FieldName, claimgroup.FieldDescription:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ClaimGroup fields.
func (cg *ClaimGroup) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case claimgroup.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			cg.ID = int(value.Int64)
		case claimgroup.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				cg.Name = value.String
			}
		case claimgroup.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				cg.Description = value.String
			}
		default:
			cg.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ClaimGroup.
// This includes values selected through modifiers, order, etc.
func (cg *ClaimGroup) Value(name string) (ent.Value, error) {
	return cg.selectValues.Get(name)
}

// QueryUsers queries the "users" edge of the ClaimGroup entity.
func (cg *ClaimGroup) QueryUsers() *UserQuery {
	return NewClaimGroupClient(cg.config).QueryUsers(cg)
}

// QueryGroupLinks queries the "group_links" edge of the ClaimGroup entity.
func (cg *ClaimGroup) QueryGroupLinks() *GroupLinkQuery {
	return NewClaimGroupClient(cg.config).QueryGroupLinks(cg)
}

// QueryClaims queries the "claims" edge of the ClaimGroup entity.
func (cg *ClaimGroup) QueryClaims() *ClaimQuery {
	return NewClaimGroupClient(cg.config).QueryClaims(cg)
}

// Update returns a builder for updating this ClaimGroup.
// Note that you need to call ClaimGroup.Unwrap() before calling this method if this ClaimGroup
// was returned from a transaction, and the transaction was committed or rolled back.
func (cg *ClaimGroup) Update() *ClaimGroupUpdateOne {
	return NewClaimGroupClient(cg.config).UpdateOne(cg)
}

// Unwrap unwraps the ClaimGroup entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (cg *ClaimGroup) Unwrap() *ClaimGroup {
	_tx, ok := cg.config.driver.(*txDriver)
	if !ok {
		panic("ent: ClaimGroup is not a transactional entity")
	}
	cg.config.driver = _tx.drv
	return cg
}

// String implements the fmt.Stringer.
func (cg *ClaimGroup) String() string {
	var builder strings.Builder
	builder.WriteString("ClaimGroup(")
	builder.WriteString(fmt.Sprintf("id=%v, ", cg.ID))
	builder.WriteString("name=")
	builder.WriteString(cg.Name)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(cg.Description)
	builder.WriteByte(')')
	return builder.String()
}

// ClaimGroups is a parsable slice of ClaimGroup.
type ClaimGroups []*ClaimGroup
