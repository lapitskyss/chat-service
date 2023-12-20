package managerhandler_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/chat-service/internal/types"
	managerhandler "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler"
)

func TestJSONEventReader_Smoke(t *testing.T) {
	chatID := types.MustParse[types.ChatID]("c920d118-f3ea-423b-b710-4c4da7610808")
	reqID := types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132")
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ManagerTypingEvent",
			"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	r := managerhandler.JSONEventReader{}

	event, err := r.Read(buf)
	require.NoError(t, err)
	assert.Equal(t, managerhandler.NewManagerTypingEvent(chatID, reqID), event)
}
