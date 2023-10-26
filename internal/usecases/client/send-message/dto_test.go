package sendmessage_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
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
				ClientID:    types.NewUserID(),
				MessageBody: "Hello, guys!",
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: sendmessage.Request{
				ID:          types.RequestIDNil,
				ClientID:    types.NewUserID(),
				MessageBody: "Hello, guys!",
			},
			wantErr: true,
		},
		{
			name: "require client id",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ClientID:    types.UserIDNil,
				MessageBody: "Hello, guys!",
			},
			wantErr: true,
		},
		{
			name: "require message body",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ClientID:    types.NewUserID(),
				MessageBody: "",
			},
			wantErr: true,
		},
		{
			name: "too big message body",
			request: sendmessage.Request{
				ID:          types.NewRequestID(),
				ClientID:    types.NewUserID(),
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
