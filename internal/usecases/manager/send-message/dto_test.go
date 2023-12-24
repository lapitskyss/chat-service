package sendmessage_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request sendmessage.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ManagerID:   types.NewUserID(),
				ChatID:      types.NewChatID(),
				MessageBody: "How can I help you, sir?",
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: sendmessage.Request{
				ID:          types.RequestIDNil,
				ManagerID:   types.NewUserID(),
				ChatID:      types.NewChatID(),
				MessageBody: "How can I help you, sir?",
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ManagerID:   types.UserIDNil,
				ChatID:      types.NewChatID(),
				MessageBody: "How can I help you, sir?",
			},
			wantErr: true,
		},
		{
			name: "require chat id",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ManagerID:   types.NewUserID(),
				ChatID:      types.ChatIDNil,
				MessageBody: "How can I help you, sir?",
			},
			wantErr: true,
		},
		{
			name: "require message body",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ManagerID:   types.NewUserID(),
				ChatID:      types.NewChatID(),
				MessageBody: "",
			},
			wantErr: true,
		},
		{
			name: "too big message body",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ManagerID:   types.NewUserID(),
				ChatID:      types.NewChatID(),
				MessageBody: strings.Repeat("x", 3001),
			},
			wantErr: true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
