package managerevents_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/chat-service/internal/types"
	managerevents "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler/events"
)

func TestDecode_DecodeJSONError(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`{"eventType": ""`))

	_, err := managerevents.Decode(buf)
	require.Error(t, err)
}

func TestDecode_IncorrectEventTypeError(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "",
		}
	`))

	_, err := managerevents.Decode(buf)
	require.Error(t, err)
}

func TestDecode_ManagerTypingEvent_Success(t *testing.T) {
	reqID := types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132")
	chatID := types.MustParse[types.ChatID]("c920d118-f3ea-423b-b710-4c4da7610808")
	buf := bytes.NewBuffer([]byte(`
		{
			"eventType": "ManagerTypingEvent",
			"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
			"requestId": "d85154f7-867b-44f1-a930-f37151568132"
		}
	`))

	res, err := managerevents.Decode(buf)
	require.NoError(t, err)
	assert.Equal(t, managerevents.ManagerTypingEvent{
		EventType: "ManagerTypingEvent",
		ChatId:    chatID,
		RequestId: reqID,
	}, res)
}
