package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/lapitskyss/chat-service/internal/types"
)

// Message holds the schema definition for the Message entity.
type Message struct {
	ent.Schema
}

// Fields of the Message.
func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.MessageID{}).
			Default(types.NewMessageID).
			Immutable(),
		//field.UUID("chat_id", types.ChatID{}),
		//field.UUID("problem_id", types.ProblemID{}),
		field.UUID("author_id", types.UserID{}).
			Optional(),
		field.Bool("is_visible_for_client").
			Default(false),
		field.Bool("is_visible_for_manager").
			Default(false),
		field.String("body").
			MaxLen(3000).
			NotEmpty(),
		field.Time("checked_at").
			Optional().
			Nillable(),
		field.Bool("is_blocked").
			Default(false),
		field.Bool("is_service").
			Default(false).
			Immutable(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("chat", Chat.Type).
			Ref("messages").
			Unique(),
		edge.From("problem", Problem.Type).
			Ref("messages").
			Unique(),
	}
}
