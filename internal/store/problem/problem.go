// Code generated by ent, DO NOT EDIT.

package problem

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/lapitskyss/chat-service/internal/types"
)

const (
	// Label holds the string label denoting the problem type in the database.
	Label = "problem"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldManagerID holds the string denoting the manager_id field in the database.
	FieldManagerID = "manager_id"
	// FieldResolvedAt holds the string denoting the resolved_at field in the database.
	FieldResolvedAt = "resolved_at"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// EdgeMessages holds the string denoting the messages edge name in mutations.
	EdgeMessages = "messages"
	// EdgeChat holds the string denoting the chat edge name in mutations.
	EdgeChat = "chat"
	// Table holds the table name of the problem in the database.
	Table = "problems"
	// MessagesTable is the table that holds the messages relation/edge.
	MessagesTable = "messages"
	// MessagesInverseTable is the table name for the Message entity.
	// It exists in this package in order to avoid circular dependency with the "message" package.
	MessagesInverseTable = "messages"
	// MessagesColumn is the table column denoting the messages relation/edge.
	MessagesColumn = "problem_messages"
	// ChatTable is the table that holds the chat relation/edge.
	ChatTable = "problems"
	// ChatInverseTable is the table name for the Chat entity.
	// It exists in this package in order to avoid circular dependency with the "chat" package.
	ChatInverseTable = "chats"
	// ChatColumn is the table column denoting the chat relation/edge.
	ChatColumn = "chat_problems"
)

// Columns holds all SQL columns for problem fields.
var Columns = []string{
	FieldID,
	FieldManagerID,
	FieldResolvedAt,
	FieldCreatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "problems"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"chat_problems",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() types.ProblemID
)

// OrderOption defines the ordering options for the Problem queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByManagerID orders the results by the manager_id field.
func ByManagerID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldManagerID, opts...).ToFunc()
}

// ByResolvedAt orders the results by the resolved_at field.
func ByResolvedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldResolvedAt, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByMessagesCount orders the results by messages count.
func ByMessagesCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newMessagesStep(), opts...)
	}
}

// ByMessages orders the results by messages terms.
func ByMessages(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newMessagesStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByChatField orders the results by chat field.
func ByChatField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newChatStep(), sql.OrderByField(field, opts...))
	}
}
func newMessagesStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(MessagesInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, MessagesTable, MessagesColumn),
	)
}
func newChatStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ChatInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, ChatTable, ChatColumn),
	)
}
