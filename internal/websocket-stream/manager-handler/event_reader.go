package managerhandler

import (
	"fmt"
	"io"

	managerevents "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler/events"
)

// EventReader converts data to event stream.
type EventReader interface {
	Read(r io.Reader) (Event, error)
}

type JSONEventReader struct{}

func (JSONEventReader) Read(r io.Reader) (Event, error) {
	e, err := managerevents.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode manager event, %v", err)
	}
	switch t := e.(type) {
	case managerevents.ManagerTypingEvent:
		return NewManagerTypingEvent(t.ChatId, t.RequestId), nil
	default:
		return nil, fmt.Errorf("unexpeced event type, %v", err)
	}
}
