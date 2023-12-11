package errhandler_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lapitskyss/chat-service/internal/server-manager/errhandler"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
)

func TestResponseBuilder(t *testing.T) {
	t.Run("with details", func(t *testing.T) {
		err := errhandler.ResponseBuilder(1000, "hello", "world")

		resp, ok := err.(errhandler.Response)
		require.True(t, ok)
		require.IsType(t, managerv1.Error{}, resp.Error)

		assert.Equal(t, managerv1.ErrorCode(1000), resp.Error.Code)
		assert.Equal(t, "hello", resp.Error.Message)
		require.NotNil(t, resp.Error.Details)
		assert.Equal(t, "world", *resp.Error.Details)
	})

	t.Run("without details", func(t *testing.T) {
		err := errhandler.ResponseBuilder(1001, "hello", "")

		resp, ok := err.(errhandler.Response)
		require.True(t, ok)
		require.IsType(t, managerv1.Error{}, resp.Error)

		assert.Equal(t, managerv1.ErrorCode(1001), resp.Error.Code)
		assert.Equal(t, "hello", resp.Error.Message)
		assert.Nil(t, resp.Error.Details)
	})
}
