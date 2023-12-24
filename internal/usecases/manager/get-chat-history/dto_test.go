package getchathistory_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/lapitskyss/chat-service/internal/types"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
)

func TestRequest_Validate(t *testing.T) {
	cases := []struct {
		name    string
		request getchathistory.Request
		wantErr bool
	}{
		// Positive.
		{
			name: "cursor specified",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  0,
				Cursor:    "eyJwYWdlX3NpemUiOjUwLCJsYXN0IjoxNjcwNTAyNTAyfQ==", // {"page_size":50,"last":1670502502}
			},
			wantErr: false,
		},
		{
			name: "page size specified",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  50,
				Cursor:    "",
			},
			wantErr: false,
		},

		// Negative.
		{
			name: "neither cursor nor pagesize specified",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  0,
				Cursor:    "",
			},
			wantErr: true,
		},
		{
			name: "cursor and pagesize specified",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  50,
				Cursor:    "eyJwYWdlX3NpemUiOjUwLCJsYXN0IjoxNjcwNTAyNTAyfQ==", // {"page_size":50,"last":1670502502}
			},
			wantErr: true,
		},
		{
			name: "require request id",
			request: getchathistory.Request{
				ID:        types.RequestIDNil,
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  10,
				Cursor:    "",
			},
			wantErr: true,
		},
		{
			name: "require manager id",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.UserIDNil,
				ChatID:    types.NewChatID(),
				PageSize:  10,
				Cursor:    "",
			},
			wantErr: true,
		},
		{
			name: "require chat id",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.ChatIDNil,
				PageSize:  10,
				Cursor:    "",
			},
			wantErr: true,
		},
		{
			name: "invalid cursor encoding",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  0,
				Cursor:    "eyJwYWdlX3NpemUiOjUwLCJs YXN0IjoxNjcwNTAyNTAyfQ", // With space.
			},
			wantErr: true,
		},
		{
			name: "too small page size",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  9,
				Cursor:    "",
			},
			wantErr: true,
		},
		{
			name: "too big page size",
			request: getchathistory.Request{
				ID:        types.NewRequestID(),
				ManagerID: types.NewUserID(),
				ChatID:    types.NewChatID(),
				PageSize:  101,
				Cursor:    "",
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
