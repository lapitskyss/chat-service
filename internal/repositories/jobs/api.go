package jobsrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
)

var ErrNoJobs = errors.New("no jobs found")

type Job struct {
	ID       types.JobID
	Name     string
	Payload  string
	Attempts int
}

func (r *Repo) FindAndReserveJob(ctx context.Context, until time.Time) (Job, error) {
	query := `
	WITH sq AS (
		SELECT id FROM jobs
		WHERE available_at <= now() AND reserved_until <= now()
		LIMIT 1 FOR UPDATE SKIP LOCKED
	)
	UPDATE jobs AS j
	SET attempts = attempts + 1, reserved_until = $1 
	FROM sq 
	WHERE sq.id = j.id RETURNING
		j.id,
		j.name,
		j.payload,
		j.attempts;
	`

	rows, err := r.db.Job(ctx).QueryContext(ctx, query, until)
	if err != nil {
		return Job{}, fmt.Errorf("query context: %w", err)
	}
	defer rows.Close()

	var j Job
	for rows.Next() {
		err = rows.Scan(&j.ID, &j.Name, &j.Payload, &j.Attempts)
		if err != nil {
			return Job{}, fmt.Errorf("scan rows: %v", err)
		}
	}

	if err = rows.Err(); err != nil {
		return Job{}, fmt.Errorf("rows err: %v", err)
	}

	if j.ID.IsZero() {
		return Job{}, ErrNoJobs
	}

	return j, nil
}

func (r *Repo) CreateJob(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	j, err := r.db.Job(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetAvailableAt(availableAt).
		Save(ctx)
	if err != nil {
		return types.JobIDNil, fmt.Errorf("create job: %v", err)
	}
	return j.ID, nil
}

func (r *Repo) CreateFailedJob(ctx context.Context, name, payload, reason string) error {
	err := r.db.FailedJob(ctx).
		Create().
		SetName(name).
		SetPayload(payload).
		SetReason(reason).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("create failed job: %v", err)
	}
	return nil
}

func (r *Repo) DeleteJob(ctx context.Context, jobID types.JobID) error {
	err := r.db.Job(ctx).
		DeleteOneID(jobID).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("delete job: %v", err)
	}
	return nil
}
