package clientevents_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clientevents "github.com/lapitskyss/chat-service/internal/server-client/events"
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
			ev: eventstream.NewMessageSentEvent(
				types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
				types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
			),
			expJSON: `{
				"eventId": "d0ffbd36-bc30-11ed-8286-461e464ebed8",
				"eventType": "MessageSentEvent",
				"messageId": "cb36a888-bc30-11ed-b843-461e464ebed8",
				"requestId": "cee5f290-bc30-11ed-b7fe-461e464ebed8"
			}`,
		},

		{
			name: "service message",
			ev: eventstream.NewNewMessageEvent(
				types.MustParse[types.EventID]("d0ffbd36-bc30-11ed-8286-461e464ebed8"),
				types.MustParse[types.RequestID]("cee5f290-bc30-11ed-b7fe-461e464ebed8"),
				types.MustParse[types.ChatID]("31b4dc06-bc31-11ed-93cc-461e464ebed8"),
				types.MustParse[types.MessageID]("cb36a888-bc30-11ed-b843-461e464ebed8"),
				types.UserIDNil,
				time.Unix(1, 1).UTC(),
				"Manager will coming soon",
				true,
			),
			expJSON: `{
				"body": "Manager will coming soon",
				"createdAt": "1970-01-01T00:00:01.000000001Z",
				"eventId": "d0ffbd36-bc30-11ed-8286-461e464ebed8",
				"eventType": "NewMessageEvent",
				"isService": true,
				"messageId": "cb36a888-bc30-11ed-b843-461e464ebed8",
				"requestId": "cee5f290-bc30-11ed-b7fe-461e464ebed8"
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
			adapted, err := clientevents.Adapter{}.Adapt(tt.ev)
			require.NoError(t, err)

			raw, err := json.Marshal(adapted)
			require.NoError(t, err)
			assert.JSONEq(t, tt.expJSON, string(raw))
		})
	}
}
