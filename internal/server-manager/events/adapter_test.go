package managerevents_test

import (
	"encoding/json"
	"testing"
	"time"

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
			name: "new chat event",
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
		{
			name: "new chat event",
			ev: eventstream.NewNewMessageEvent(
				types.MustParse[types.EventID]("8cfd1993-9a7b-45c4-9e9c-2a33086a860f"),
				types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132"),
				types.MustParse[types.ChatID]("c920d118-f3ea-423b-b710-4c4da7610808"),
				types.MustParse[types.MessageID]("66293272-4b7c-4846-a351-6dadef6317c8"),
				types.MustParse[types.UserID]("7d67b14d-221e-4499-9be2-6707d7df1adc"),
				time.Unix(1, 1).UTC(),
				"Чего там с деньгами",
				false,
			),
			expJSON: `{
				"authorId": "7d67b14d-221e-4499-9be2-6707d7df1adc",
				"body": "Чего там с деньгами",
				"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
				"createdAt": "1970-01-01T00:00:01.000000001Z",
				"eventId": "8cfd1993-9a7b-45c4-9e9c-2a33086a860f",
				"eventType": "NewMessageEvent",
				"messageId": "66293272-4b7c-4846-a351-6dadef6317c8",
				"requestId": "d85154f7-867b-44f1-a930-f37151568132"
			}`,
		},
		{
			name: "chat close event",
			ev: eventstream.NewChatClosedEvent(
				types.MustParse[types.EventID]("8cfd1993-9a7b-45c4-9e9c-2a33086a860f"),
				types.MustParse[types.ChatID]("c920d118-f3ea-423b-b710-4c4da7610808"),
				types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132"),
				false,
			),
			expJSON: `{
				"eventId": "8cfd1993-9a7b-45c4-9e9c-2a33086a860f",
				"eventType": "ChatClosedEvent",
				"chatId": "c920d118-f3ea-423b-b710-4c4da7610808",
				"requestId": "d85154f7-867b-44f1-a930-f37151568132",
				"canTakeMoreProblems": false
			}`,
		},
		{
			name: "typing event",
			ev: eventstream.NewTypingEvent(
				types.MustParse[types.EventID]("8cfd1993-9a7b-45c4-9e9c-2a33086a860f"),
				types.MustParse[types.UserID]("7d67b14d-221e-4499-9be2-6707d7df1adc"),
				types.MustParse[types.RequestID]("d85154f7-867b-44f1-a930-f37151568132"),
			),
			expJSON: `{
				"eventId": "8cfd1993-9a7b-45c4-9e9c-2a33086a860f",
				"clientId": "7d67b14d-221e-4499-9be2-6707d7df1adc",
				"eventType": "TypingEvent",
				"requestId": "d85154f7-867b-44f1-a930-f37151568132"
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
