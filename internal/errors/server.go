package errors

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// ServerError is used to return custom error codes to client.
type ServerError struct {
	Code    int
	Message string
	cause   error
}

func NewServerError(code int, msg string, err error) *ServerError {
	return &ServerError{
		Code:    code,
		Message: msg,
		cause:   err,
	}
}

func (s *ServerError) Error() string {
	return fmt.Sprintf("%s: %v", s.Message, s.cause)
}

func (s *ServerError) Is(target error) bool {
	return errors.Is(s.cause, target)
}

func GetServerErrorCode(err error) int {
	code, _, _ := ProcessServerError(err)
	return code
}

// ProcessServerError tries to retrieve from given error it's code, message and some details.
// For example, that fields can be used to build error response for client.
func ProcessServerError(err error) (code int, msg string, details string) {
	var svcErr *ServerError
	if errors.As(err, &svcErr) {
		return svcErr.Code, svcErr.Message, svcErr.Error()
	}
	var httpErr *echo.HTTPError
	if errors.As(err, &httpErr) {
		return httpErr.Code, fmt.Sprintf("%v", httpErr.Message), httpErr.Error()
	}
	return 500, "something went wrong", err.Error()
}
