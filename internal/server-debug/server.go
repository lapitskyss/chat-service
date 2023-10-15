package serverdebug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/buildinfo"
	"github.com/lapitskyss/chat-service/internal/logger"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
)

const (
	readHeaderTimeout = time.Second
	shutdownTimeout   = 3 * time.Second
)

//go:generate options-gen -out-filename=server_options.gen.go -from-struct=Options
type Options struct {
	addr string `option:"mandatory" validate:"required,hostname_port"`
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
	e.Use(middleware.Recover())

	s := &Server{
		lg: serverLogger(),
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
	e.GET("/schema/client", s.ClientSchema)

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

func (s *Server) ClientSchema(c echo.Context) error {
	swagger, err := clientv1.GetSwagger()
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, swagger)
}
