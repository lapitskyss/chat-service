package jobsrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lapitskyss/chat-service/internal/store"
	"github.com/lapitskyss/chat-service/internal/store/job"
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
	var j *store.Job
	err := r.db.RunInTx(ctx, func(ctx context.Context) error {
		var err error
		now := time.Now()
		j, err = r.db.Job(ctx).
			Query().
			Where(
				job.ReservedUntilLTE(now),
				job.AvailableAtLTE(now),
			).
			ForUpdate().
			First(ctx)
		if err != nil {
			if store.IsNotFound(err) {
				return ErrNoJobs
			}
			return fmt.Errorf("find job: %v", err)
		}

		j, err = r.db.Job(ctx).
			UpdateOne(j).
			SetReservedUntil(until).
			SetAttempts(j.Attempts + 1).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("reserve job: %v", err)
		}

		return nil
	})
	if err != nil {
		return Job{}, fmt.Errorf("find and reserve job: %w", err)
	}
	return Job{
		ID:       j.ID,
		Name:     j.Name,
		Payload:  j.Payload,
		Attempts: j.Attempts,
	}, nil
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
