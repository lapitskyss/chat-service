package getchats_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request getchats.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: getchats.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: getchats.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: getchats.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
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
