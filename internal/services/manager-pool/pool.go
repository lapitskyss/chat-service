package managerpool

import (
	"context"
	"errors"
	"io"

	"github.com/lapitskyss/chat-service/internal/types"
)

var ErrNoAvailableManagers = errors.New("no available managers")

// Pool represents concurrent-safe FIFO queue.
type Pool interface {
	io.Closer
	Get(ctx context.Context) (types.UserID, error)
	Put(ctx context.Context, managerID types.UserID) error
	Contains(ctx context.Context, managerID types.UserID) (bool, error)
	Size() int
}
