// Code generated by ent, DO NOT EDIT.

package store

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/store/message"
	"github.com/lapitskyss/chat-service/internal/store/predicate"
	"github.com/lapitskyss/chat-service/internal/store/problem"
	"github.com/lapitskyss/chat-service/internal/types"
)

// ChatUpdate is the builder for updating Chat entities.
type ChatUpdate struct {
	config
	hooks    []Hook
	mutation *ChatMutation
}

// Where appends a list predicates to the ChatUpdate builder.
func (cu *ChatUpdate) Where(ps ...predicate.Chat) *ChatUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// AddMessageIDs adds the "messages" edge to the Message entity by IDs.
func (cu *ChatUpdate) AddMessageIDs(ids ...types.MessageID) *ChatUpdate {
	cu.mutation.AddMessageIDs(ids...)
	return cu
}

// AddMessages adds the "messages" edges to the Message entity.
func (cu *ChatUpdate) AddMessages(m ...*Message) *ChatUpdate {
	ids := make([]types.MessageID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cu.AddMessageIDs(ids...)
}

// AddProblemIDs adds the "problems" edge to the Problem entity by IDs.
func (cu *ChatUpdate) AddProblemIDs(ids ...types.ProblemID) *ChatUpdate {
	cu.mutation.AddProblemIDs(ids...)
	return cu
}

// AddProblems adds the "problems" edges to the Problem entity.
func (cu *ChatUpdate) AddProblems(p ...*Problem) *ChatUpdate {
	ids := make([]types.ProblemID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cu.AddProblemIDs(ids...)
}

// Mutation returns the ChatMutation object of the builder.
func (cu *ChatUpdate) Mutation() *ChatMutation {
	return cu.mutation
}

// ClearMessages clears all "messages" edges to the Message entity.
func (cu *ChatUpdate) ClearMessages() *ChatUpdate {
	cu.mutation.ClearMessages()
	return cu
}

// RemoveMessageIDs removes the "messages" edge to Message entities by IDs.
func (cu *ChatUpdate) RemoveMessageIDs(ids ...types.MessageID) *ChatUpdate {
	cu.mutation.RemoveMessageIDs(ids...)
	return cu
}

// RemoveMessages removes "messages" edges to Message entities.
func (cu *ChatUpdate) RemoveMessages(m ...*Message) *ChatUpdate {
	ids := make([]types.MessageID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cu.RemoveMessageIDs(ids...)
}

// ClearProblems clears all "problems" edges to the Problem entity.
func (cu *ChatUpdate) ClearProblems() *ChatUpdate {
	cu.mutation.ClearProblems()
	return cu
}

// RemoveProblemIDs removes the "problems" edge to Problem entities by IDs.
func (cu *ChatUpdate) RemoveProblemIDs(ids ...types.ProblemID) *ChatUpdate {
	cu.mutation.RemoveProblemIDs(ids...)
	return cu
}

// RemoveProblems removes "problems" edges to Problem entities.
func (cu *ChatUpdate) RemoveProblems(p ...*Problem) *ChatUpdate {
	ids := make([]types.ProblemID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cu.RemoveProblemIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *ChatUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *ChatUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *ChatUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *ChatUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *ChatUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(chat.Table, chat.Columns, sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if cu.mutation.MessagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedMessagesIDs(); len(nodes) > 0 && !cu.mutation.MessagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.MessagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cu.mutation.ProblemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.RemovedProblemsIDs(); len(nodes) > 0 && !cu.mutation.ProblemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cu.mutation.ProblemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chat.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// ChatUpdateOne is the builder for updating a single Chat entity.
type ChatUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *ChatMutation
}

// AddMessageIDs adds the "messages" edge to the Message entity by IDs.
func (cuo *ChatUpdateOne) AddMessageIDs(ids ...types.MessageID) *ChatUpdateOne {
	cuo.mutation.AddMessageIDs(ids...)
	return cuo
}

// AddMessages adds the "messages" edges to the Message entity.
func (cuo *ChatUpdateOne) AddMessages(m ...*Message) *ChatUpdateOne {
	ids := make([]types.MessageID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cuo.AddMessageIDs(ids...)
}

// AddProblemIDs adds the "problems" edge to the Problem entity by IDs.
func (cuo *ChatUpdateOne) AddProblemIDs(ids ...types.ProblemID) *ChatUpdateOne {
	cuo.mutation.AddProblemIDs(ids...)
	return cuo
}

// AddProblems adds the "problems" edges to the Problem entity.
func (cuo *ChatUpdateOne) AddProblems(p ...*Problem) *ChatUpdateOne {
	ids := make([]types.ProblemID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cuo.AddProblemIDs(ids...)
}

// Mutation returns the ChatMutation object of the builder.
func (cuo *ChatUpdateOne) Mutation() *ChatMutation {
	return cuo.mutation
}

// ClearMessages clears all "messages" edges to the Message entity.
func (cuo *ChatUpdateOne) ClearMessages() *ChatUpdateOne {
	cuo.mutation.ClearMessages()
	return cuo
}

// RemoveMessageIDs removes the "messages" edge to Message entities by IDs.
func (cuo *ChatUpdateOne) RemoveMessageIDs(ids ...types.MessageID) *ChatUpdateOne {
	cuo.mutation.RemoveMessageIDs(ids...)
	return cuo
}

// RemoveMessages removes "messages" edges to Message entities.
func (cuo *ChatUpdateOne) RemoveMessages(m ...*Message) *ChatUpdateOne {
	ids := make([]types.MessageID, len(m))
	for i := range m {
		ids[i] = m[i].ID
	}
	return cuo.RemoveMessageIDs(ids...)
}

// ClearProblems clears all "problems" edges to the Problem entity.
func (cuo *ChatUpdateOne) ClearProblems() *ChatUpdateOne {
	cuo.mutation.ClearProblems()
	return cuo
}

// RemoveProblemIDs removes the "problems" edge to Problem entities by IDs.
func (cuo *ChatUpdateOne) RemoveProblemIDs(ids ...types.ProblemID) *ChatUpdateOne {
	cuo.mutation.RemoveProblemIDs(ids...)
	return cuo
}

// RemoveProblems removes "problems" edges to Problem entities.
func (cuo *ChatUpdateOne) RemoveProblems(p ...*Problem) *ChatUpdateOne {
	ids := make([]types.ProblemID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return cuo.RemoveProblemIDs(ids...)
}

// Where appends a list predicates to the ChatUpdate builder.
func (cuo *ChatUpdateOne) Where(ps ...predicate.Chat) *ChatUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *ChatUpdateOne) Select(field string, fields ...string) *ChatUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Chat entity.
func (cuo *ChatUpdateOne) Save(ctx context.Context) (*Chat, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *ChatUpdateOne) SaveX(ctx context.Context) *Chat {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *ChatUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *ChatUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *ChatUpdateOne) sqlSave(ctx context.Context) (_node *Chat, err error) {
	_spec := sqlgraph.NewUpdateSpec(chat.Table, chat.Columns, sqlgraph.NewFieldSpec(chat.FieldID, field.TypeUUID))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`store: missing "Chat.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, chat.FieldID)
		for _, f := range fields {
			if !chat.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("store: invalid field %q for query", f)}
			}
			if f != chat.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if cuo.mutation.MessagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedMessagesIDs(); len(nodes) > 0 && !cuo.mutation.MessagesCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.MessagesIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.MessagesTable,
			Columns: []string{chat.MessagesColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(message.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if cuo.mutation.ProblemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.RemovedProblemsIDs(); len(nodes) > 0 && !cuo.mutation.ProblemsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := cuo.mutation.ProblemsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   chat.ProblemsTable,
			Columns: []string{chat.ProblemsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(problem.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Chat{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{chat.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
