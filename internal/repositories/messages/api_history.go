package messagesrepo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"

	"github.com/lapitskyss/chat-service/internal/store/chat"
	"github.com/lapitskyss/chat-service/internal/store/message"
	"github.com/lapitskyss/chat-service/internal/types"
)

var (
	ErrInvalidPageSize      = errors.New("invalid page size")
	ErrInvalidCursor        = errors.New("invalid cursor")
	ErrInvalidRequestParams = errors.New("invalid request params")
)

type Cursor struct {
	LastCreatedAt time.Time
	PageSize      int
}

// GetClientChatMessages returns Nth page of messages in the chat for client side.
func (r *Repo) GetClientChatMessages(
	ctx context.Context,
	clientID types.UserID,
	pageSize int,
	cursor *Cursor,
) ([]Message, *Cursor, error) {
	pageSize, createdFrom, err := validateGetClientChatMessages(pageSize, cursor)
	if err != nil {
		return nil, nil, fmt.Errorf("validate request: %w", err)
	}

	query := r.db.Message(ctx).
		Query().
		Unique(false).
		Where(message.HasChatWith(chat.ClientID(clientID))).
		Where(message.IsVisibleForClient(true)).
		Order(message.ByCreatedAt(sql.OrderDesc())).
		Limit(pageSize + 1)

	if !createdFrom.IsZero() {
		query = query.Where(message.CreatedAtLT(createdFrom))
	}

	msgs, err := query.All(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("get client chat messags: %v", err)
	}

	if len(msgs) <= pageSize {
		return adaptStoreMessages(msgs), nil, nil
	}

	c := &Cursor{
		LastCreatedAt: msgs[len(msgs)-2].CreatedAt,
		PageSize:      pageSize,
	}
	return adaptStoreMessages(msgs[:len(msgs)-1]), c, nil
}

func validateGetClientChatMessages(pageSize int, cursor *Cursor) (int, time.Time, error) {
	if cursor != nil {
		if !isPageSizeValid(cursor.PageSize) {
			return 0, time.Time{}, ErrInvalidCursor
		}
		if cursor.LastCreatedAt.IsZero() {
			return 0, time.Time{}, ErrInvalidCursor
		}
		return cursor.PageSize, cursor.LastCreatedAt, nil
	}

	if pageSize != 0 {
		if !isPageSizeValid(pageSize) {
			return 0, time.Time{}, ErrInvalidPageSize
		}
		return pageSize, time.Time{}, nil
	}

	return 0, time.Time{}, ErrInvalidRequestParams
}

func isPageSizeValid(pageSize int) bool {
	return pageSize >= 10 && pageSize <= 100
}
