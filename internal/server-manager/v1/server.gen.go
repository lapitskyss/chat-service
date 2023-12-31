// Package managerv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package managerv1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"github.com/lapitskyss/chat-service/internal/types"
	"github.com/oapi-codegen/runtime"
)

const (
	BearerAuthScopes = "bearerAuth.Scopes"
)

// Defines values for ErrorCode.
const (
	ErrorCodeManagerOverloaded     ErrorCode = 5000
	ErrorCodeNoActiveProblemInChat ErrorCode = 5001
)

// Chat defines model for Chat.
type Chat struct {
	ChatId   types.ChatID `json:"chatId"`
	ClientId types.UserID `json:"clientId"`
}

// ChatId defines model for ChatId.
type ChatId struct {
	ChatId types.ChatID `json:"chatId"`
}

// ChatList defines model for ChatList.
type ChatList struct {
	Chats []Chat `json:"chats"`
}

// CloseChat defines model for CloseChat.
type CloseChat = interface{}

// CloseChatRequest defines model for CloseChatRequest.
type CloseChatRequest struct {
	ChatId types.ChatID `json:"chatId"`
}

// CloseChatResponse defines model for CloseChatResponse.
type CloseChatResponse struct {
	Data  *CloseChat `json:"data,omitempty"`
	Error *Error     `json:"error,omitempty"`
}

// Error defines model for Error.
type Error struct {
	// Code contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
	Code    ErrorCode `json:"code"`
	Details *string   `json:"details,omitempty"`
	Message string    `json:"message"`
}

// ErrorCode contains HTTP error codes and specific business logic error codes (the last must be >= 1000).
type ErrorCode int

// FreeHands defines model for FreeHands.
type FreeHands = interface{}

// FreeHandsBtnAvailability defines model for FreeHandsBtnAvailability.
type FreeHandsBtnAvailability struct {
	Available bool `json:"available"`
}

// FreeHandsResponse defines model for FreeHandsResponse.
type FreeHandsResponse struct {
	Data  *FreeHands `json:"data,omitempty"`
	Error *Error     `json:"error,omitempty"`
}

// GetChatHistoryRequest defines model for GetChatHistoryRequest.
type GetChatHistoryRequest struct {
	ChatId   types.ChatID `json:"chatId"`
	Cursor   *string      `json:"cursor,omitempty"`
	PageSize *int         `json:"pageSize,omitempty"`
}

// GetChatHistoryResponse defines model for GetChatHistoryResponse.
type GetChatHistoryResponse struct {
	Data  *MessagesPage `json:"data,omitempty"`
	Error *Error        `json:"error,omitempty"`
}

// GetChatsResponse defines model for GetChatsResponse.
type GetChatsResponse struct {
	Data  *ChatList `json:"data,omitempty"`
	Error *Error    `json:"error,omitempty"`
}

// GetFreeHandsBtnAvailabilityResponse defines model for GetFreeHandsBtnAvailabilityResponse.
type GetFreeHandsBtnAvailabilityResponse struct {
	Data  *FreeHandsBtnAvailability `json:"data,omitempty"`
	Error *Error                    `json:"error,omitempty"`
}

// Message defines model for Message.
type Message struct {
	AuthorId  types.UserID    `json:"authorId"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessageWithoutBody defines model for MessageWithoutBody.
type MessageWithoutBody struct {
	AuthorId  types.UserID    `json:"authorId"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
	Next     string    `json:"next"`
}

// SendMessageRequest defines model for SendMessageRequest.
type SendMessageRequest struct {
	ChatId      types.ChatID `json:"chatId"`
	MessageBody string       `json:"messageBody"`
}

// SendMessageResponse defines model for SendMessageResponse.
type SendMessageResponse struct {
	Data  *MessageWithoutBody `json:"data,omitempty"`
	Error *Error              `json:"error,omitempty"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostCloseChatParams defines parameters for PostCloseChat.
type PostCloseChatParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostFreeHandsParams defines parameters for PostFreeHands.
type PostFreeHandsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatHistoryParams defines parameters for PostGetChatHistory.
type PostGetChatHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetChatsParams defines parameters for PostGetChats.
type PostGetChatsParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetFreeHandsBtnAvailabilityParams defines parameters for PostGetFreeHandsBtnAvailability.
type PostGetFreeHandsBtnAvailabilityParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostSendMessageParams defines parameters for PostSendMessage.
type PostSendMessageParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostCloseChatJSONRequestBody defines body for PostCloseChat for application/json ContentType.
type PostCloseChatJSONRequestBody = CloseChatRequest

// PostGetChatHistoryJSONRequestBody defines body for PostGetChatHistory for application/json ContentType.
type PostGetChatHistoryJSONRequestBody = GetChatHistoryRequest

// PostSendMessageJSONRequestBody defines body for PostSendMessage for application/json ContentType.
type PostSendMessageJSONRequestBody = SendMessageRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /closeChat)
	PostCloseChat(ctx echo.Context, params PostCloseChatParams) error

	// (POST /freeHands)
	PostFreeHands(ctx echo.Context, params PostFreeHandsParams) error

	// (POST /getChatHistory)
	PostGetChatHistory(ctx echo.Context, params PostGetChatHistoryParams) error

	// (POST /getChats)
	PostGetChats(ctx echo.Context, params PostGetChatsParams) error

	// (POST /getFreeHandsBtnAvailability)
	PostGetFreeHandsBtnAvailability(ctx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error

	// (POST /sendMessage)
	PostSendMessage(ctx echo.Context, params PostSendMessageParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostCloseChat converts echo context to params.
func (w *ServerInterfaceWrapper) PostCloseChat(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostCloseChatParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostCloseChat(ctx, params)
	return err
}

// PostFreeHands converts echo context to params.
func (w *ServerInterfaceWrapper) PostFreeHands(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostFreeHandsParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostFreeHands(ctx, params)
	return err
}

// PostGetChatHistory converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetChatHistory(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetChatHistoryParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostGetChatHistory(ctx, params)
	return err
}

// PostGetChats converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetChats(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetChatsParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostGetChats(ctx, params)
	return err
}

// PostGetFreeHandsBtnAvailability converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetFreeHandsBtnAvailability(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetFreeHandsBtnAvailabilityParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostGetFreeHandsBtnAvailability(ctx, params)
	return err
}

// PostSendMessage converts echo context to params.
func (w *ServerInterfaceWrapper) PostSendMessage(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostSendMessageParams

	headers := ctx.Request().Header
	// ------------- Required header parameter "X-Request-ID" -------------
	if valueList, found := headers[http.CanonicalHeaderKey("X-Request-ID")]; found {
		var XRequestID XRequestIDHeader
		n := len(valueList)
		if n != 1 {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Expected one value for X-Request-ID, got %d", n))
		}

		err = runtime.BindStyledParameterWithLocation("simple", false, "X-Request-ID", runtime.ParamLocationHeader, valueList[0], &XRequestID)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter X-Request-ID: %s", err))
		}

		params.XRequestID = XRequestID
	} else {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Header parameter X-Request-ID is required, but not found"))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostSendMessage(ctx, params)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/closeChat", wrapper.PostCloseChat)
	router.POST(baseURL+"/freeHands", wrapper.PostFreeHands)
	router.POST(baseURL+"/getChatHistory", wrapper.PostGetChatHistory)
	router.POST(baseURL+"/getChats", wrapper.PostGetChats)
	router.POST(baseURL+"/getFreeHandsBtnAvailability", wrapper.PostGetFreeHandsBtnAvailability)
	router.POST(baseURL+"/sendMessage", wrapper.PostSendMessage)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RYW2/bRhP9K4v5vocWoEyqboBAQB+cq10kjVGnSABXDytyJG5C7jK7Q8Wqwf9e7HJ5",
	"kUjZju0YDvokkXubOXNm9gwvIVZ5oSRKMjC7hIJrniOhdk8f/8QvJRo6eXGMPEFt3wkJM0jrxwAkzxFm",
	"8HHiZ05OXkAAGr+UQmMCM9IlBmDiFHNuVy+VzjnBDMpSJBAAbQq73pAWcgUBXExWauJf2h9z0JrQH52I",
	"vFCaaosphRmsBKXl4iBWeZjxQpD5vDEmjFNOE4N6LWIMhSTUkmeh2xiqqqoa05y3z1PuduRZ9m4Js/NL",
	"+L/GJczgf2EHUugXhHb2SQJVcAmFVgVqEui2iTOB0g7dyt2/DOrv4ms/KuedkfPWKLX4hDFBNa8C8M7N",
	"Br6177/dM7fnA3hWm9h48UYYGvfD/RGEuftzXaihan3kWvMNjJ1r6mMzZdCTqf/omfzDo9r5YwolDQ4d",
	"Sjjxa0FtUaoCQK2Vvm7FSzfJmfaymb8DpErwRrs8txOrABIkLjK3dhvqKoAcjeErHBnbwaaZGNTnzxv7",
	"nntrEjSxFgUJZUtnrCRxIQ07fv/+lDnHmV1nGJcJMwXGYilitiiNkGgMy9RKxFvzfqIUWcYNsbw0xBbI",
	"/i6j6BB/Y9Moin4+gABQljnMzp9EURQ8iaLpPIBcSJHbt79GUUssG/aVK+QXE7tmsubalnRj/WqdeMsl",
	"X6F+t0adKZ6gZWY7+Ic6ikms8VSrRYb5iXQxtSC80ojHXCamzoP28RnJozUXGV+ITNBmGEZej2Z97BdK",
	"ZcjlAPxu7taRd+NmZ/ktuPkayUJwLAwpvfmhkj6AuNSm9naQDwVf4Zn4x0Ga84uaTFNLppZa0yGzrigk",
	"u0DdJWJv6xw0pzYRbx+0O/KmvXBuZ8G+DLknMu/m3S2MfNvVxJ2ULSlV+rGJngAWKtmMsjnWyAmTI9qy",
	"OOGEExI5DsyuAhC39M6D9gC3tDOoDYX3vu/rvIvhB0GpKumZB+jHCOd/Imqj4aor2yBQXnvcXMs2GTyQ",
	"swFIvKAbqx0DfoG18Qxl4jfu3Xd37KP8QQ1Bc37xBuXKgn4Y+VuneTENbma022u849ly4R5uon56fXOh",
	"tX0pxqUWtDmzY/XpC+Qa9VFpPW6eXjW8/v3De/DdrJNLbrQjekpU1MwTcqlclAVZgQXPuPzMzsrC8prZ",
	"YDAv99jR6QkEsEZtauW6nlpPVIGSFwJmcHgQHRxC4DLBGRjGvc4HClXTYFv+WpyZxK/Mh4SRYlbP2oyx",
	"ytXCze1cW3zgVBnqGoVg6/PEHl51U8LB54tqXtMCTVv3rBxHWfO1KDIRu8PDT8Zae9n7cnGjXqZh/06O",
	"ky7RvaiJ5cD6JYq+x/meus6AbeQ9MZkvLweeaOGyJ9P3BK1hhEaebNhSaRfB+hvCaMw6AX1PMftOwA37",
	"hRHgvEXMrs+wD91qS8Hux+81Ess9hpboLK1XjBN+Wxc/XtaPNzoPTP09TcR+/huWCUO7ITRXB8913MIQ",
	"U0sXQMO+CkqZrYWsqBtfc2UwH3siDPqfEQDdhAF6V7b1ewEVhtm6w1K7ki1KIiVZ18vvQXLvWY8e3Gtb",
	"uxG8jxo4WJxi/Hmk/phOtNznjdvTQo+3+oxozgcuPWOS8ab3bk/hOVT72u58bjGz3UOD+faGL3CNmSpy",
	"lMTqWRBAqTMv82ZhmKmYZ6kyNHsaPZ2GVrjNq38DAAD//yHhsRnrGQAA",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
