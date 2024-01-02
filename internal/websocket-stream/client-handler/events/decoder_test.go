package clientevents_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/chat-service/internal/types"
	clientevents "github.com/lapitskyss/chat-service/internal/websocket-stream/client-handler/events"
)

func TestDecode_DecodeJSONError(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`{"eventType": ""`))

	_, err := clientevents.Decode(buf)
	require.Error(t, err)
}

func TestDecode_IncorrectEventTypeError(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "",
		}
	`))

	_, err := clientevents.Decode(buf)
	require.Error(t, err)
}

func TestDecode_ClientTypingEvent_Success(t *testing.T) {
	reqID := types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132")
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ClientTypingEvent",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	res, err := clientevents.Decode(buf)
	require.NoError(t, err)
	assert.Equal(t, clientevents.ClientTypingEvent{
		EventType: "ClientTypingEvent",
		RequestId: reqID,
	}, res)
}
