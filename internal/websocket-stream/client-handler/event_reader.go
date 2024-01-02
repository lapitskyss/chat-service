package clienthandler

import (
	"fmt"
	"io"

	clientevents "github.com/lapitskyss/chat-service/internal/websocket-stream/client-handler/events"
)

// EventReader converts data to event stream.
type EventReader interface {
	Read(r io.Reader) (Event, error)
}

type JSONEventReader struct{}

func (JSONEventReader) Read(r io.Reader) (Event, error) {
	e, err := clientevents.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decode client event, %v", err)
	}
	switch t := e.(type) {
	case clientevents.ClientTypingEvent:
		return NewClientTypingEvent(t.RequestId), nil
	default:
		return nil, fmt.Errorf("unexpeced event type, %v", err)
	}
}
