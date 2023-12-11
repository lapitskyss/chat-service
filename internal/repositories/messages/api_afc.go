package messagesrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/lapitskyss/chat-service/internal/types"
)

func (r *Repo) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	err := r.db.Message(ctx).
		UpdateOneID(msgID).
		SetIsVisibleForManager(true).
		SetCheckedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("mark as visible for manager: %v", err)
	}
	return nil
}

func (r *Repo) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	err := r.db.Message(ctx).
		UpdateOneID(msgID).
		SetIsBlocked(true).
		SetCheckedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("block messsage: %v", err)
	}
	return nil
}
