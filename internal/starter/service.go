package starter

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/sync/errgroup"

	"github.com/lapitskyss/chat-service/internal/server"
	serverdebug "github.com/lapitskyss/chat-service/internal/server-debug"
	afcverdictsprocessor "github.com/lapitskyss/chat-service/internal/services/afc-verdicts-processor"
	managerscheduler "github.com/lapitskyss/chat-service/internal/services/manager-scheduler"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
)

type Service struct {
	ctx context.Context

	srvClient  *server.Server
	srvManager *server.Server
	srvDebug   *serverdebug.Server

	outBox              *outbox.Service
	managerScheduler    *managerscheduler.Service
	afcVerdictProcessor *afcverdictsprocessor.Service
}

func (s *Service) Run(cleanup func()) error {
	defer cleanup()

	eg, ctx := errgroup.WithContext(s.ctx)

	// Run servers.
	eg.Go(func() error { return s.srvClient.Run(ctx) })
	eg.Go(func() error { return s.srvManager.Run(ctx) })
	eg.Go(func() error { return s.srvDebug.Run(ctx) })

	// Run services
	eg.Go(func() error { return s.outBox.Run(ctx) })
	eg.Go(func() error { return s.managerScheduler.Run(ctx) })
	eg.Go(func() error { return s.afcVerdictProcessor.Run(ctx) })

	err := eg.Wait()
	if err != nil && !errors.Is(err, context.Canceled) {
		return fmt.Errorf("wait app stop: %v", err)
	}
	return nil
}

func NewService(
	ctx context.Context,
	srvClient ServerClient,
	srvManager ServerManager,
	srvDebug *serverdebug.Server,
	outBox *outbox.Service,
	managerScheduler *managerscheduler.Service,
	afcVerdictProcessor *afcverdictsprocessor.Service,
) *Service {
	return &Service{
		ctx:                 ctx,
		srvClient:           srvClient,
		srvManager:          srvManager,
		srvDebug:            srvDebug,
		outBox:              outBox,
		managerScheduler:    managerScheduler,
		afcVerdictProcessor: afcVerdictProcessor,
	}
}
