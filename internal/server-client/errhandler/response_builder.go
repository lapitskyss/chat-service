package errhandler

import (
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	"github.com/lapitskyss/chat-service/pkg/pointer"
)

type Response struct {
	Error clientv1.Error `json:"error"`
}

var ResponseBuilder = func(code int, msg string, details string) any {
	return Response{
		Error: clientv1.Error{
			Code:    clientv1.ErrorCode(code),
			Message: msg,
			Details: pointer.PtrWithZeroAsNil[string](details),
		},
	}
}
