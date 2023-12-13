package managerevents_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	managerevents "github.com/lapitskyss/chat-service/internal/server-manager/events"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/types"
)

func TestAdapter_Adapt(t *testing.T) {
	cases := []struct {
		name    string
		ev      eventstream.Event
		expJSON string
	}{
		{
			name: "smoke",
			ev: eventstream.NewNewChatEvent(
				types.MustParse[types.EventID]("9d55940f-751a-4608-8dd3-5a3ae04598a0"),
				types.MustParse[types.ChatID]("5e121000-ce14-4b47-bcbd-ee82f6343039"),
				types.MustParse[types.UserID]("7d67b14d-221e-4499-9be2-6707d7df1adc"),
				types.MustParse[types.RequestID]("88b37e1e-78fa-4329-b66d-60eb955a0ff4"),
				true,
			),
			expJSON: `{
				"canTakeMoreProblems": true,
				"chatId": "5e121000-ce14-4b47-bcbd-ee82f6343039",
				"clientId": "7d67b14d-221e-4499-9be2-6707d7df1adc",
				"eventId": "9d55940f-751a-4608-8dd3-5a3ae04598a0",
				"eventType": "NewChatEvent",
				"requestId": "88b37e1e-78fa-4329-b66d-60eb955a0ff4"
			}`,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			adapted, err := managerevents.Adapter{}.Adapt(tt.ev)
			require.NoError(t, err)

			raw, err := json.Marshal(adapted)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expJSON, string(raw))
		})
	}
}
