package canreceiveproblems_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request canreceiveproblems.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "valid request",
			request: canreceiveproblems.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "require request id",
			request: canreceiveproblems.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: canreceiveproblems.Request{
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
