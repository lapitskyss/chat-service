package clienthandler_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/chat-service/internal/types"
	clienthandler "github.com/lapitskyss/chat-service/internal/websocket-stream/client-handler"
)

func TestJSONEventReader_Smoke(t *testing.T) {
	reqID := types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132")
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ClientTypingEvent",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	r := clienthandler.JSONEventReader{}

	event, err := r.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, clienthandler.NewClientTypingEvent(reqID), event)
}
