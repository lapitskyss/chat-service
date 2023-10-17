package serverclient

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	oapimdlwr "github.com/deepmap/oapi-codegen/pkg/middleware"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/middlewares"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	logger               *zap.Logger              `option:"mandatory" validate:"required"`
	addr                 string                   `option:"mandatory" validate:"required,hostname_port"`
	allowOrigins         []string                 `option:"mandatory" validate:"min=1"`
	v1Swagger            *openapi3.T              `option:"mandatory" validate:"required"`
	v1Handlers           clientv1.ServerInterface `option:"mandatory" validate:"required"`
	keycloakIntrospector middlewares.Introspector `option:"mandatory" validate:"required"`
	authResource         string                   `option:"mandatory" validate:"required"`
	authRole             string                   `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	e := echo.New()
	e.Use(middlewares.Logger(opts.logger))
	e.Use(middlewares.Recover(opts.logger))
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: opts.allowOrigins,
		AllowMethods: []string{http.MethodPost},
	}))
	e.Use(middleware.BodyLimit("12K")) // 3000 characters * 4 байт
	e.Use(middlewares.NewKeycloakTokenAuth(opts.keycloakIntrospector, opts.authResource, opts.authRole))

	v1 := e.Group("v1", oapimdlwr.OapiRequestValidatorWithOptions(opts.v1Swagger, &oapimdlwr.Options{
		Options: openapi3filter.Options{
			ExcludeRequestBody:  false,
			ExcludeResponseBody: true,
			AuthenticationFunc:  openapi3filter.NoopAuthenticationFunc,
		},
		SilenceServersWarning: true,
	}))
	clientv1.RegisterHandlers(v1, opts.v1Handlers)

	return &Server{
		lg: opts.logger,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}, nil
}

func (s *Server) Run(ctx context.Context) error {
	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return s.srv.Shutdown(ctx)
	})

	eg.Go(func() error {
		s.lg.Info("listen and serve", zap.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %v", err)
		}
		return nil
	})

	return eg.Wait()
}
