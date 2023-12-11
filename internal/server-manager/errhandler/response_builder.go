package errhandler

import (
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

type Response struct {
	Error managerv1.Error `json:"error"`
}

func ResponseBuilder(code int, msg string, details string) any {
	return Response{
		Error: managerv1.Error{
			Code:    managerv1.ErrorCode(code),
			Message: msg,
			Details: pointer.PtrWithZeroAsNil(details),
		},
	}
}
