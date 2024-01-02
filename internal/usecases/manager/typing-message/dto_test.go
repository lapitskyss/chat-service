package managertypingmessage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	managertypingmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/typing-message"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request managertypingmessage.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: managertypingmessage.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: managertypingmessage.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: managertypingmessage.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
				ChatID:    types.NewChatID(),
			},
			wantErr: true,
		},
		{
			name: "require chat id",
			request: managertypingmessage.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.ChatIDNil,
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
