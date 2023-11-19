package serverclient

import (
	"fmt"
	"net/http"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	"github.com/lapitskyss/chat-service/internal/server"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger           *zap.Logger              `option:"mandatory" validate:"required"`
	addr             string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins     []string                 `option:"mandatory" validate:"min=1"`
	introspector     middlewares.Introspector `option:"mandatory" validate:"required"`
	requiredResource string                   `option:"mandatory" validate:"required"`
	requiredRole     string                   `option:"mandatory" validate:"required"`
	v1Swagger        *openapi3.T              `option:"mandatory" validate:"required"`
	v1Handlers       clientv1.ServerInterface `option:"mandatory" validate:"required"`
	httpErrorHandler echo.HTTPErrorHandler    `option:"mandatory" validate:"required"`
}

func New(opts Options) (*server.Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	e := echo.New()
	e.HTTPErrorHandler = opts.httpErrorHandler
	e.Use(middlewares.Logger(opts.logger))
	e.Use(middlewares.Recover(opts.logger))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: opts.allowOrigins,
		AllowMethods: []string{http.MethodPost},
	}))
	e.Use(middleware.BodyLimit("12K")) // 3000 characters * 4 байт
	e.Use(middlewares.NewKeycloakTokenAuth(opts.introspector, opts.requiredResource, opts.requiredRole))

	v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
		SilenceServersWarning: true,
	}))
	clientv1.RegisterHandlers(v1, opts.v1Handlers)

	return server.New(server.NewOptions(opts.logger, opts.addr, e))
}
