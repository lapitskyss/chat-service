package managerevents

import (
	"errors"
	"fmt"

	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	websocketstream "github.com/lapitskyss/chat-service/internal/websocket-stream"
)

var ErrUnexpectedEventType = errors.New("unexpected event type")

var _ websocketstream.EventAdapter = Adapter{}

type Adapter struct{}

func (Adapter) Adapt(ev eventstream.Event) (any, error) {
	switch e := ev.(type) {
	case *eventstream.NewChatEvent:
		event := Event{}

		err := event.FromNewChatEvent(NewChatEvent{
			CanTakeMoreProblems: e.CanTakeMoreProblems,
			ChatId:              e.ChatID,
			ClientId:            e.ClientID,
			EventId:             e.EventID,
			RequestId:           e.RequestID,
		})
		if err != nil {
			return nil, fmt.Errorf("from new message event: %v", err)
		}

		return event, nil
	case *eventstream.NewMessageEvent:
		event := Event{}

		err := event.FromNewMessageEvent(NewMessageEvent{
			AuthorId:  e.UserID,
			Body:      e.MessageBody,
			ChatId:    e.ChatID,
			CreatedAt: e.CreatedAt,
			EventId:   e.EventID,
			MessageId: e.MessageID,
			RequestId: e.RequestID,
		})
		if err != nil {
			return nil, fmt.Errorf("from new message event: %v", err)
		}

		return event, nil
	default:
		return nil, ErrUnexpectedEventType
	}
}
