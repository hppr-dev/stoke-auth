// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"stoke/internal/ent/claimgroup"
	"stoke/internal/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// ClaimGroupDelete is the builder for deleting a ClaimGroup entity.
type ClaimGroupDelete struct {
	config
	hooks    []Hook
	mutation *ClaimGroupMutation
}

// Where appends a list predicates to the ClaimGroupDelete builder.
func (cgd *ClaimGroupDelete) Where(ps ...predicate.ClaimGroup) *ClaimGroupDelete {
	cgd.mutation.Where(ps...)
	return cgd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cgd *ClaimGroupDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cgd.sqlExec, cgd.mutation, cgd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cgd *ClaimGroupDelete) ExecX(ctx context.Context) int {
	n, err := cgd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cgd *ClaimGroupDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(claimgroup.Table, sqlgraph.NewFieldSpec(claimgroup.FieldID, field.TypeInt))
	if ps := cgd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cgd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cgd.mutation.done = true
	return affected, err
}

// ClaimGroupDeleteOne is the builder for deleting a single ClaimGroup entity.
type ClaimGroupDeleteOne struct {
	cgd *ClaimGroupDelete
}

// Where appends a list predicates to the ClaimGroupDelete builder.
func (cgdo *ClaimGroupDeleteOne) Where(ps ...predicate.ClaimGroup) *ClaimGroupDeleteOne {
	cgdo.cgd.mutation.Where(ps...)
	return cgdo
}

// Exec executes the deletion query.
func (cgdo *ClaimGroupDeleteOne) Exec(ctx context.Context) error {
	n, err := cgdo.cgd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{claimgroup.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cgdo *ClaimGroupDeleteOne) ExecX(ctx context.Context) {
	if err := cgdo.Exec(ctx); err != nil {
		panic(err)
	}
}
