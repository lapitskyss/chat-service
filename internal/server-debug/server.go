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

	e.GET("/", index.handler)
	e.GET("/version", s.Version)
	e.GET("/debug/pprof/", s.Pprof)
	e.GET("/debug/pprof/profile", s.PprofProfile)
	e.PUT("/log/level", s.ChangeLogLevel)

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

func (s *Server) Pprof(c echo.Context) error {
	pprof.Index(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) PprofProfile(c echo.Context) error {
	pprof.Profile(c.Response().Writer, c.Request())
	return nil
}

func (s *Server) ChangeLogLevel(c echo.Context) error {
	level := c.FormValue("level")
	err := logger.Init(logger.NewOptions(level))
	if err != nil {
		return err
	}
	s.lg = serverLogger()
	return nil
}
