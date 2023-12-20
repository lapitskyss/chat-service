package clienttypingmessage_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	clienttypingmessage "github.com/lapitskyss/chat-service/internal/usecases/client/typing-message"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request clienttypingmessage.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: clienttypingmessage.Request{
				ID:       types.NewRequestID(),
				ClientID: types.NewUserID(),
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: clienttypingmessage.Request{
				ID:       types.RequestIDNil,
				ClientID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "require client id",
			request: clienttypingmessage.Request{
				ID:       types.NewRequestID(),
				ClientID: types.UserIDNil,
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
