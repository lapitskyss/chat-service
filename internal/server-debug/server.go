package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/buildinfo"
	"github.com/lapitskyss/chat-service/internal/logger"
	"github.com/lapitskyss/chat-service/internal/middlewares"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr                 string      `option:"mandatory" validate:"required,hostname_port"`
	clientSwagger        *openapi3.T `option:"mandatory" validate:"required"`
	managerSwagger       *openapi3.T `option:"mandatory" validate:"required"`
	clientEventsSwagger  *openapi3.T `option:"mandatory" validate:"required"`
	managerEventsSwagger *openapi3.T `option:"mandatory" validate:"required"`
}

type Server struct {
	lg  *zap.Logger
	srv *http.Server
}

func New(opts Options) (*Server, error) {
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("validate options: %v", err)
	}

	lg := serverLogger()

	e := echo.New()
	e.Use(middlewares.Recover(lg))
	e.Use(middlewares.Logger(lg))

	s := &Server{
		lg: lg,
		srv: &http.Server{
			Addr:              opts.addr,
			Handler:           e,
			ReadHeaderTimeout: readHeaderTimeout,
		},
	}

	index := newIndexPage()
	index.addPage("/version", "Get build information")
	index.addPage("/debug/pprof/", "Go std profiler")
	index.addPage("/debug/pprof/profile?seconds=30", "Take half-min profile")
	index.addPage("/debug/error", "Debug Sentry error event")
	index.addPage("/schema/client", "Get client OpenAPI specification")
	index.addPage("/schema/manager", "Get manager OpenAPI specification")
	index.addPage("/schema/clientEvents", "Get client events OpenAPI specification")

	e.GET("/", index.handler)
	e.GET("/version", s.Version)
	{
		pprofMux := http.NewServeMux()
		pprofMux.HandleFunc("/debug/pprof/", pprof.Index)
		pprofMux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		pprofMux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		pprofMux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		pprofMux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		e.GET("/debug/pprof/*", echo.WrapHandler(pprofMux))
	}
	e.PUT("/log/level", echo.WrapHandler(logger.Level))
	e.GET("/debug/error", s.DebugError)
	e.GET("/schema/client", s.ExposeSchema(opts.clientSwagger))
	e.GET("/schema/manager", s.ExposeSchema(opts.managerSwagger))
	e.GET("/schema/clientEvents", s.ExposeSchema(opts.clientEventsSwagger))
	e.GET("/schema/managerEvents", s.ExposeSchema(opts.managerEventsSwagger))

	return s, nil
}

func serverLogger() *zap.Logger {
	return zap.L().Named("server-debug")
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

func (s *Server) Version(c echo.Context) error {
	return c.JSON(http.StatusOK, buildinfo.BuildInfo)
}

func (s *Server) DebugError(c echo.Context) error {
	s.lg.Error("look for me in the Sentry")
	return c.String(http.StatusOK, "event send")
}

func (s *Server) ExposeSchema(swagger *openapi3.T) echo.HandlerFunc {
	return func(eCtx echo.Context) error {
		return eCtx.JSON(http.StatusOK, swagger)
	}
}
