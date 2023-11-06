// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lapitskyss/chat-service/internal/store/failedjob"
	"github.com/lapitskyss/chat-service/internal/store/predicate"
)

// FailedJobUpdate is the builder for updating FailedJob entities.
type FailedJobUpdate struct {
	config
	hooks    []Hook
	mutation *FailedJobMutation
}

// Where appends a list predicates to the FailedJobUpdate builder.
func (fju *FailedJobUpdate) Where(ps ...predicate.FailedJob) *FailedJobUpdate {
	fju.mutation.Where(ps...)
	return fju
}

// Mutation returns the FailedJobMutation object of the builder.
func (fju *FailedJobUpdate) Mutation() *FailedJobMutation {
	return fju.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fju *FailedJobUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, fju.sqlSave, fju.mutation, fju.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fju *FailedJobUpdate) SaveX(ctx context.Context) int {
	affected, err := fju.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fju *FailedJobUpdate) Exec(ctx context.Context) error {
	_, err := fju.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fju *FailedJobUpdate) ExecX(ctx context.Context) {
	if err := fju.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fju *FailedJobUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(failedjob.Table, failedjob.Columns, sqlgraph.NewFieldSpec(failedjob.FieldID, field.TypeUUID))
	if ps := fju.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fju.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{failedjob.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fju.mutation.done = true
	return n, nil
}

// FailedJobUpdateOne is the builder for updating a single FailedJob entity.
type FailedJobUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FailedJobMutation
}

// Mutation returns the FailedJobMutation object of the builder.
func (fjuo *FailedJobUpdateOne) Mutation() *FailedJobMutation {
	return fjuo.mutation
}

// Where appends a list predicates to the FailedJobUpdate builder.
func (fjuo *FailedJobUpdateOne) Where(ps ...predicate.FailedJob) *FailedJobUpdateOne {
	fjuo.mutation.Where(ps...)
	return fjuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fjuo *FailedJobUpdateOne) Select(field string, fields ...string) *FailedJobUpdateOne {
	fjuo.fields = append([]string{field}, fields...)
	return fjuo
}

// Save executes the query and returns the updated FailedJob entity.
func (fjuo *FailedJobUpdateOne) Save(ctx context.Context) (*FailedJob, error) {
	return withHooks(ctx, fjuo.sqlSave, fjuo.mutation, fjuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fjuo *FailedJobUpdateOne) SaveX(ctx context.Context) *FailedJob {
	node, err := fjuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fjuo *FailedJobUpdateOne) Exec(ctx context.Context) error {
	_, err := fjuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fjuo *FailedJobUpdateOne) ExecX(ctx context.Context) {
	if err := fjuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (fjuo *FailedJobUpdateOne) sqlSave(ctx context.Context) (_node *FailedJob, err error) {
	_spec := sqlgraph.NewUpdateSpec(failedjob.Table, failedjob.Columns, sqlgraph.NewFieldSpec(failedjob.FieldID, field.TypeUUID))
	id, ok := fjuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`store: missing "FailedJob.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fjuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, failedjob.FieldID)
		for _, f := range fields {
			if !failedjob.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("store: invalid field %q for query", f)}
			}
			if f != failedjob.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fjuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	_node = &FailedJob{config: fjuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fjuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{failedjob.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fjuo.mutation.done = true
	return _node, nil
}
