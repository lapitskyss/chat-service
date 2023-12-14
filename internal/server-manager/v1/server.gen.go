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
	ErrorCodeManagerOverloaded ErrorCode = 5000
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

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostFreeHandsParams defines parameters for PostFreeHands.
type PostFreeHandsParams struct {
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

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /freeHands)
	PostFreeHands(ctx echo.Context, params PostFreeHandsParams) error

	// (POST /getChats)
	PostGetChats(ctx echo.Context, params PostGetChatsParams) error

	// (POST /getFreeHandsBtnAvailability)
	PostGetFreeHandsBtnAvailability(ctx echo.Context, params PostGetFreeHandsBtnAvailabilityParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
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

	router.POST(baseURL+"/freeHands", wrapper.PostFreeHands)
	router.POST(baseURL+"/getChats", wrapper.PostGetChats)
	router.POST(baseURL+"/getFreeHandsBtnAvailability", wrapper.PostGetFreeHandsBtnAvailability)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/8xW32/bNhD+V4jbHjZAtpRlAwoBe0iTtcmwoUGbYQUyP9DS2WJDkSx5cmsE+t+HoyRb",
	"ju157X6gT7bII+/7vrv7pEcobO2sQUMB8kdw0ssaCX18evsa3zcY6ObqGmWJnteUgRyq7jEBI2uEHN5O",
	"+sjJzRUk4PF9ozyWkJNvMIFQVFhLPr2wvpYEOTSNKiEBWjs+H8grs4QEPk6WdtIv8k+YbiCMdyeqdtZT",
	"h5gqyGGpqGrm08LWqZZOUXhYh5AWlaRJQL9SBabKEHojdRovhrZt2wFaZHtZyXij1PrVAvL7R/ja4wJy",
	"+CrdipT2B1KOvimhTR7BeevQk8J4TaEVGt76LLq/BfT/CddxVe63IGcbUHb+DguCdtYm0JPL97ht1j+d",
	"Wbzzf2DWQRxY/KICHeYR/yjCOv45VWpoNxyl93INh/KGmPYn760/kNOWeCpTPHrJgW0CJZJUOp7dVbdN",
	"oMYQ5BIP7D2BNQQmXf4NvsseTYmh8MqRsjzXhTUklQni+u7uViAHCj4XhDSlCA4LtVCFmDdBGQxBaLtU",
	"xU7cN1Sh0DKQqJtAYo7ijybLzvFHcZZl2bdTSABNU0N+/0OWZbMEamVUzQvfZ9lGYi7yMhrMxwmHT1bS",
	"s9UEprTB/6s0con+1Qq9trLEruovPOK1NCULN358TuZiJZWWc6UVrfcLJLtdPVZ1bq1GafZk3cbupHyN",
	"wVkTcP/yUpI8Vf0t8jYBHLroZL90M/ASifv0H0LYjMznITgm9r+ky9MSfjJINnwsGq9o/Yb3OgxzlB79",
	"RcMWNDy9GPzt59/voH9NxIaIu1vDq4hcR1+ZhY2do4hbCJ5L8yDeNI4tTrCuou9XcXF7Awms0Idu6lZn",
	"zMQ6NNIpyOF8mk3PIYmmGAGmi1FPg7Odpe2O7nC5R1muxcJ6YfCD6Hwe4u1ecijbN9zasK1VzLR97R95",
	"721D0r3PgnbG09GVOEL8Lss6yzPE6Xm4nNOqiAjSd4ERP44+C/5W9TdNFOXeZd8jEnxeI2E57cudLvu5",
	"OK7cSyQRXUsFEnYhopOLD4oqwTURztu5xjpMD8o4zN0XruKePRwQMQZEGcbq/aWBHhVUBcFNKyo+KeYN",
	"kTVi65pHlDya64sX96TzHdD7YpBDFBUWD0+bd2RWkfLYpu5nTIi/iQZBdq++whVq62o0JLooSKDxunes",
	"PE21LaSubKD8WfbsLGUPmrV/BgAA///Fksi9DwwAAA==",
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
