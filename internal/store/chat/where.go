// Code generated by ent, DO NOT EDIT.

package chat

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/lapitskyss/chat-service/internal/store/predicate"
	"github.com/lapitskyss/chat-service/internal/types"
)

// ID filters vertices based on their ID field.
func ID(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id types.ChatID) predicate.Chat {
	return predicate.Chat(sql.FieldLTE(FieldID, id))
}

// ClientID applies equality check predicate on the "client_id" field. It's identical to ClientIDEQ.
func ClientID(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldClientID, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldCreatedAt, v))
}

// ClientIDEQ applies the EQ predicate on the "client_id" field.
func ClientIDEQ(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldClientID, v))
}

// ClientIDNEQ applies the NEQ predicate on the "client_id" field.
func ClientIDNEQ(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldNEQ(FieldClientID, v))
}

// ClientIDIn applies the In predicate on the "client_id" field.
func ClientIDIn(vs ...types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldIn(FieldClientID, vs...))
}

// ClientIDNotIn applies the NotIn predicate on the "client_id" field.
func ClientIDNotIn(vs ...types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldNotIn(FieldClientID, vs...))
}

// ClientIDGT applies the GT predicate on the "client_id" field.
func ClientIDGT(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldGT(FieldClientID, v))
}

// ClientIDGTE applies the GTE predicate on the "client_id" field.
func ClientIDGTE(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldGTE(FieldClientID, v))
}

// ClientIDLT applies the LT predicate on the "client_id" field.
func ClientIDLT(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldLT(FieldClientID, v))
}

// ClientIDLTE applies the LTE predicate on the "client_id" field.
func ClientIDLTE(v types.UserID) predicate.Chat {
	return predicate.Chat(sql.FieldLTE(FieldClientID, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.Chat {
	return predicate.Chat(sql.FieldLTE(FieldCreatedAt, v))
}

// HasMessages applies the HasEdge predicate on the "messages" edge.
func HasMessages() predicate.Chat {
	return predicate.Chat(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, MessagesTable, MessagesColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasMessagesWith applies the HasEdge predicate on the "messages" edge with a given conditions (other predicates).
func HasMessagesWith(preds ...predicate.Message) predicate.Chat {
	return predicate.Chat(func(s *sql.Selector) {
		step := newMessagesStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProblems applies the HasEdge predicate on the "problems" edge.
func HasProblems() predicate.Chat {
	return predicate.Chat(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, ProblemsTable, ProblemsColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProblemsWith applies the HasEdge predicate on the "problems" edge with a given conditions (other predicates).
func HasProblemsWith(preds ...predicate.Problem) predicate.Chat {
	return predicate.Chat(func(s *sql.Selector) {
		step := newProblemsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.Chat) predicate.Chat {
	return predicate.Chat(sql.AndPredicates(predicates...))
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.Chat) predicate.Chat {
	return predicate.Chat(sql.OrPredicates(predicates...))
}

// Not applies the not operator on the given predicate.
func Not(p predicate.Chat) predicate.Chat {
	return predicate.Chat(sql.NotPredicates(p))
}
