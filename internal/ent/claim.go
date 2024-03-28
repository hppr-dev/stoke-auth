// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"stoke/internal/ent/claim"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
)

// Claim is the model entity for the Claim schema.
type Claim struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// ShortName holds the value of the "short_name" field.
	ShortName string `json:"short_name,omitempty"`
	// Value holds the value of the "value" field.
	Value string `json:"value,omitempty"`
	// Description holds the value of the "description" field.
	Description string `json:"description,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ClaimQuery when eager-loading is set.
	Edges        ClaimEdges `json:"edges"`
	selectValues sql.SelectValues
}

// ClaimEdges holds the relations/edges for other nodes in the graph.
type ClaimEdges struct {
	// ClaimGroups holds the value of the claim_groups edge.
	ClaimGroups []*ClaimGroup `json:"claim_groups,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ClaimGroupsOrErr returns the ClaimGroups value or an error if the edge
// was not loaded in eager-loading.
func (e ClaimEdges) ClaimGroupsOrErr() ([]*ClaimGroup, error) {
	if e.loadedTypes[0] {
		return e.ClaimGroups, nil
	}
	return nil, &NotLoadedError{edge: "claim_groups"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Claim) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case claim.FieldID:
			values[i] = new(sql.NullInt64)
		case claim.FieldName, claim.FieldShortName, claim.FieldValue, claim.FieldDescription:
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Claim fields.
func (c *Claim) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case claim.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			c.ID = int(value.Int64)
		case claim.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				c.Name = value.String
			}
		case claim.FieldShortName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field short_name", values[i])
			} else if value.Valid {
				c.ShortName = value.String
			}
		case claim.FieldValue:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field value", values[i])
			} else if value.Valid {
				c.Value = value.String
			}
		case claim.FieldDescription:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field description", values[i])
			} else if value.Valid {
				c.Description = value.String
			}
		default:
			c.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// GetValue returns the ent.Value that was dynamically selected and assigned to the Claim.
// This includes values selected through modifiers, order, etc.
func (c *Claim) GetValue(name string) (ent.Value, error) {
	return c.selectValues.Get(name)
}

// QueryClaimGroups queries the "claim_groups" edge of the Claim entity.
func (c *Claim) QueryClaimGroups() *ClaimGroupQuery {
	return NewClaimClient(c.config).QueryClaimGroups(c)
}

// Update returns a builder for updating this Claim.
// Note that you need to call Claim.Unwrap() before calling this method if this Claim
// was returned from a transaction, and the transaction was committed or rolled back.
func (c *Claim) Update() *ClaimUpdateOne {
	return NewClaimClient(c.config).UpdateOne(c)
}

// Unwrap unwraps the Claim entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (c *Claim) Unwrap() *Claim {
	_tx, ok := c.config.driver.(*txDriver)
	if !ok {
		panic("ent: Claim is not a transactional entity")
	}
	c.config.driver = _tx.drv
	return c
}

// String implements the fmt.Stringer.
func (c *Claim) String() string {
	var builder strings.Builder
	builder.WriteString("Claim(")
	builder.WriteString(fmt.Sprintf("id=%v, ", c.ID))
	builder.WriteString("name=")
	builder.WriteString(c.Name)
	builder.WriteString(", ")
	builder.WriteString("short_name=")
	builder.WriteString(c.ShortName)
	builder.WriteString(", ")
	builder.WriteString("value=")
	builder.WriteString(c.Value)
	builder.WriteString(", ")
	builder.WriteString("description=")
	builder.WriteString(c.Description)
	builder.WriteByte(')')
	return builder.String()
}

// Claims is a parsable slice of Claim.
type Claims []*Claim
