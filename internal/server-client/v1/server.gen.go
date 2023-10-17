// Package clientv1 provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package clientv1

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

// GetHistoryRequest defines model for GetHistoryRequest.
type GetHistoryRequest struct {
	Cursor   *string `json:"cursor,omitempty"`
	PageSize *int    `json:"pageSize,omitempty"`
}

// GetHistoryResponse defines model for GetHistoryResponse.
type GetHistoryResponse struct {
	Data MessagesPage `json:"data"`
}

// Message defines model for Message.
type Message struct {
	AuthorId  types.UserID    `json:"authorId"`
	Body      string          `json:"body"`
	CreatedAt time.Time       `json:"createdAt"`
	Id        types.MessageID `json:"id"`
}

// MessagesPage defines model for MessagesPage.
type MessagesPage struct {
	Messages []Message `json:"messages"`
}

// XRequestIDHeader defines model for XRequestIDHeader.
type XRequestIDHeader = types.RequestID

// PostGetHistoryParams defines parameters for PostGetHistory.
type PostGetHistoryParams struct {
	XRequestID XRequestIDHeader `json:"X-Request-ID"`
}

// PostGetHistoryJSONRequestBody defines body for PostGetHistory for application/json ContentType.
type PostGetHistoryJSONRequestBody = GetHistoryRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (POST /getHistory)
	PostGetHistory(ctx echo.Context, params PostGetHistoryParams) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// PostGetHistory converts echo context to params.
func (w *ServerInterfaceWrapper) PostGetHistory(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{})

	// Parameter object where we will unmarshal all parameters from the context
	var params PostGetHistoryParams

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
	err = w.Handler.PostGetHistory(ctx, params)
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

	router.POST(baseURL+"/getHistory", wrapper.PostGetHistory)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7RUW0/cSgz+K5HPecxuwuEF5Y2LDmylSqi0KhLahyExyZTMBY+zYovy36u57I2gClXl",
	"cWLH/vx9n/0CtVHWaNTsoHoBK0goZKTwuv2CTwM6XlxcoWiQ/DepoYIuPnPQQiFUcDtLmbPFBeRA+DRI",
	"wgYqpgFzcHWHSvi/HwwpwVDBMMgGcuC19f87JqlbyOF51pqZVNYQRzjcQQWt5G64n9dGFb2wkt3j2rmi",
	"7gTPHNJK1lhIzUha9IWv6GBMpVL98HG+nQbGcdygCoNeIl9Jx4bWKSc0J2ORWGJIqQdyJjBwiHnMwYoW",
	"b+RP9EElnqUaFFRHZZmDknrz2s7qkbZIEcN+Y2eNdjjt3AgO3P1L+AAV/FPsJCvSDMVndE606K5Fi+AL",
	"7yS4iwWWYw4pa9pCDNwZWjTvluiA128OKei+Df0VCccc7k2zfpPymlAwNqd8gLgRjDOWCiewxxzkH06X",
	"SPuIAV/pFABtpUjT78+6p2FUeiKkStGwqYzKvdM4nqA0tSAS64mHtoWXYXWwHkjy+sZXid3uURDS6eDJ",
	"2Lz+35D96ftXSAvnW8Tojv2O2UY6pH4wQXDJvY+cCf2Y3QzWk52dd4Kz816i5uz0egE5rJCcNP4irY78",
	"CMaiFlZCBcfzcn4MeVAn4Cva7aoF2kxc8gZdTdJyrHKJnHnJsi5mziHUJOHjfj3g2jjeLW1osLuYd29z",
	"vUspJhd1XEae0fFZMnttNKMO6IS1vaxD9+KH8xBf9o7p73SdXrRXdvOXOXyIZydw9F9ZfgiAdNkCgkPC",
	"N27Oeul47jP27RUY3TfW3dLz5fdpw/dhuQtcYW+s8g6JWZDDQH3yWFUUvalF3xnH1Ul5UhbeNsvxVwAA",
	"AP//oL+s7AYHAAA=",
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
