package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/lapitskyss/chat-service/internal/types"
)

const messageBodyMaxLength = 3000

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).Default(types.NewMessageID).Unique().Immutable(),
		field.UUID("chat_id", types.ChatID{}),
		field.UUID("problem_id", types.ProblemID{}),
		field.UUID("author_id", types.UserID{}).Optional().Immutable(),
		field.UUID("initial_request_id", types.RequestID{}).Unique().Immutable(),
		field.Bool("is_visible_for_client").Default(false),
		field.Bool("is_visible_for_manager").Default(false),
		field.Text("body").NotEmpty().MaxLen(messageBodyMaxLength).Immutable(),
		field.Time("checked_at").Optional(),
		field.Bool("is_blocked").Default(false),
		field.Bool("is_service").Default(false).Immutable(),
		newCreatedAtField(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		// The message has one chat.
		edge.From("chat", Chat.Type).
			Ref("messages").
			Field("chat_id").
			Required().Unique(),

		// The message has one problem.
		edge.From("problem", Problem.Type).
			Ref("messages").
			Field("problem_id").
			Required().Unique(),
	}
}

func (Message) Indexes() []ent.Index {
	return []ent.Index{
		// Getting history is based on created_at field.
		index.Fields("created_at").
			Annotations(
				entsql.DescColumns("created_at"),
				entsql.IndexType("BTREE"),
			),
	}
}
