package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/server"
	servermanager "github.com/lapitskyss/chat-service/internal/server-manager"
	managererrhandler "github.com/lapitskyss/chat-service/internal/server-manager/errhandler"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/server/errhandler"
	managerload "github.com/lapitskyss/chat-service/internal/services/manager-load"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
	websocketstream "github.com/lapitskyss/chat-service/internal/websocket-stream"
)

const nameServerManager = "server-manager"

func initServerManager(
	productionMode bool,

	addr string,
	allowOrigins []string,
	secWsProtocol string,
	v1Swagger *openapi3.T,

	keycloak *keycloakclient.Client,
	requiredResource string,
	requiredRole string,

	managerLoadSvc *managerload.Service,
	managerPool managerpool.Pool,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	canReceiveProblemUserCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(
		managerLoadSvc,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("canreceiveproblems usecase: %v", err)
	}
	freeHandsUserCase, err := freehands.New(freehands.NewOptions(
		managerLoadSvc,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("canreceiveproblems usecase: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(
		canReceiveProblemUserCase,
		freeHandsUserCase,
	))
	if err != nil {
		return nil, fmt.Errorf("create v1 manager handlers: %v", err)
	}

	shutdownCh := make(chan struct{})
	cancelFn := func() {
		close(shutdownCh)
	}

	wsUpgrader := websocketstream.NewUpgrader(allowOrigins, secWsProtocol)
	wsHandler, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		lg,
		dummyEventStream{},
		dummyAdapter{},
		websocketstream.JSONEventWriter{},
		wsUpgrader,
		shutdownCh,
	))
	if err != nil {
		return nil, fmt.Errorf("websock etstream handler: %v", err)
	}

	httpErrHandler, err := errhandler.New(errhandler.NewOptions(lg, productionMode, managererrhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create errhandler: %v", err)
	}

	srv, err := servermanager.New(servermanager.NewOptions(
		lg,
		addr,
		allowOrigins,
		keycloak,
		requiredResource,
		requiredRole,
		v1Swagger,
		v1Handlers,
		wsHandler,
		httpErrHandler.Handle,
		cancelFn,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
