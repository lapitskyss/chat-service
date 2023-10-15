package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"github.com/lapitskyss/chat-service/internal/types"
)

// Chat holds the schema definition for the Chat entity.
type Chat struct {
	ent.Schema
}

// Fields of the Chat.
func (Chat) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.ChatID{}).
			Default(types.NewChatID).
			Immutable(),
		field.UUID("client_id", types.UserID{}).
			Unique(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
	}
}

// Edges of the Chat.
func (Chat) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("problems", Problem.Type),
		edge.To("messages", Message.Type),
	}
}