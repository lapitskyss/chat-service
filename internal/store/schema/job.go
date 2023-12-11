package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"

	"github.com/lapitskyss/chat-service/internal/types"
)

// jobMaxAttempts is some limit as protection from endless retries of outbox jobs.
const jobMaxAttempts = 30

type Job struct {
	ent.Schema
}

func (Job) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.JobID{}).Default(types.NewJobID).Unique().Immutable(),

		field.Text("name").
			Comment("Job name. Name determines handler.").
			NotEmpty().Immutable(),

		field.Text("payload").
			Comment("Required data to complete the job.").
			Immutable(),

		field.Int("attempts").
			Comment(`The number of execution attempts.
If a certain threshold is exceeded, the task can be removed from the queue.`).
			Min(0).Max(jobMaxAttempts).Default(0),

		field.Time("available_at").
			Comment("The time when the job becomes available for execution. Useful for delayed execution.").
			Default(time.Now).Immutable(),

		field.Time("reserved_until").
			Comment(`Until this time the task is "reserved". Used to synchronize goroutines processing the queue.
When grabbing a task, the goroutine puts in reserved_until <time.Now() + some timeout>.
Until that time the task is considered "reserved", other goroutines will skip it.`).
			Default(time.Now),

		newCreatedAtField(),
	}
}

func (Job) Indexes() []ent.Index {
	return []ent.Index{
		// Getting job to execute is based on available_at and reserved_until fields.
		index.Fields("available_at", "reserved_until"),
	}
}

type FailedJob struct {
	ent.Schema
}

func (FailedJob) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", types.FailedJobID{}).Default(types.NewFailedJobID).Unique().Immutable(),
		field.Text("name").NotEmpty().Immutable(),
		field.Text("payload").NotEmpty().Immutable(),
		field.Text("reason").NotEmpty().Immutable(),
		newCreatedAtField(),
	}
}
