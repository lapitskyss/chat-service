package sendclientmessagejob_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	msgproducer "github.com/lapitskyss/chat-service/internal/services/msg-producer"
	sendclientmessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-client-message"
	sendclientmessagejobmocks "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-client-message/mocks"
	"github.com/lapitskyss/chat-service/internal/types"
)

func TestJob_Handle(t *testing.T) {
	// Arrange.
	ctx := context.Background()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	msgProducer := sendclientmessagejobmocks.NewMockmessageProducer(ctrl)
	msgRepo := sendclientmessagejobmocks.NewMockmessageRepository(ctrl)
	eventStream := sendclientmessagejobmocks.NewMockeventStream(ctrl)
	job, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(msgProducer, msgRepo, eventStream))
	require.NoError(t, err)

	clientID := types.NewUserID()
	msgID := types.NewMessageID()
	chatID := types.NewChatID()
	const body = "Hello!"

	msg := messagesrepo.Message{
		ID:                  msgID,
		ChatID:              chatID,
		AuthorID:            clientID,
		Body:                body,
		CreatedAt:           time.Now(),
		IsVisibleForClient:  true,
		IsVisibleForManager: false,
		IsBlocked:           false,
		IsService:           false,
	}
	msgRepo.EXPECT().GetMessageByID(gomock.Any(), msgID).Return(&msg, nil)

	msgProducer.EXPECT().ProduceMessage(gomock.Any(), msgproducer.Message{
		ID:         msgID,
		ChatID:     chatID,
		Body:       body,
		FromClient: true,
	}).Return(nil)

	eventStream.EXPECT().Publish(gomock.Any(), clientID, EqPublishEventParams(msg)).
		Return(nil)

	// Action & assert.
	payload, err := sendclientmessagejob.MarshalPayload(msgID)
	require.NoError(t, err)

	err = job.Handle(ctx, payload)
	require.NoError(t, err)
}

type eqPublishEventMatcher struct {
	msg messagesrepo.Message
}

func (e eqPublishEventMatcher) Matches(x any) bool {
	arg, ok := x.(*eventstream.NewMessageEvent)
	if !ok {
		return false
	}

	if e.msg.RequestID != arg.RequestID {
		return false
	}
	if e.msg.ChatID != arg.ChatID {
		return false
	}
	if e.msg.ID != arg.MessageID {
		return false
	}
	if e.msg.AuthorID != arg.UserID {
		return false
	}
	if e.msg.CreatedAt != arg.CreatedAt {
		return false
	}
	if e.msg.Body != arg.MessageBody {
		return false
	}
	if e.msg.IsService != arg.IsService {
		return false
	}

	return true
}

func (e eqPublishEventMatcher) String() string {
	return fmt.Sprintf("matches msg %v", e.msg)
}

func EqPublishEventParams(msg messagesrepo.Message) gomock.Matcher {
	return eqPublishEventMatcher{msg}
}
