package clienttypingmessage

import (
	"context"
	"errors"
	"fmt"

	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

//go:generate mockgen -source=$GOFILE -destination=mocks/usecase_mock.gen.go -package=clienttypingmessagemocks

var ErrInvalidRequest = errors.New("invalid request")

type problemsRepository interface {
	GetClientOpenProblem(ctx context.Context, clientID types.UserID) (*problemsrepo.Problem, error)
}

type eventStream interface {
	Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error
}

//go:generate options-gen -out-filename=usecase_options.gen.go -from-struct=Options
type Options struct {
	problemRepo problemsRepository `option:"mandatory" validate:"required"`
	eventStream eventStream        `option:"mandatory" validate:"required"`
}

type UseCase struct {
	Options
}

func New(opts Options) (UseCase, error) {
	if err := opts.Validate(); err != nil {
		return UseCase{}, fmt.Errorf("validate options: %v", err)
	}
	return UseCase{Options: opts}, nil
}

func (u UseCase) Handle(ctx context.Context, req Request) error {
	if err := req.Validate(); err != nil {
		return fmt.Errorf("validate request: %w: %v", ErrInvalidRequest, err)
	}

	problem, err := u.problemRepo.GetClientOpenProblem(ctx, req.ClientID)
	if err != nil {
		return fmt.Errorf("problem repo, get chat open problem, %v", err)
	}

	// Send event NewTypingEvent to manager
	err = u.eventStream.Publish(ctx, problem.ManagerID, eventstream.NewTypingEvent(
		types.NewEventID(),
		req.ClientID,
		req.ID,
	))
	if err != nil {
		return fmt.Errorf("event stream, publish message blick event: %v", err)
	}

	return nil
}
