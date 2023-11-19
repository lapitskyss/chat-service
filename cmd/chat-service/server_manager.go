package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/server"
	servermanager "github.com/lapitskyss/chat-service/internal/server-manager"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/server/errhandler"
	managerload "github.com/lapitskyss/chat-service/internal/services/manager-load"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
)

const nameServerManager = "server-manager"

func initServerManager(
	productionMode bool,

	addr string,
	allowOrigins []string,
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

	httpErrHandler, err := errhandler.New(errhandler.NewOptions(lg, productionMode, errhandler.ResponseBuilder))
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
		httpErrHandler.Handle,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
