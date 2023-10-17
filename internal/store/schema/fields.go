package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

func newCreatedAtField() ent.Field {
	return field.Time("created_at").
		Default(time.Now).
		Immutable()
}
