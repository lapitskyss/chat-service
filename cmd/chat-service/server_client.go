package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	serverclient "github.com/lapitskyss/chat-service/internal/server-client"
	"github.com/lapitskyss/chat-service/internal/server-client/errhandler"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
)

const nameServerClient = "server-client"

func initServerClient(
	productionMode bool,

	addr string,
	allowOrigins []string,
	v1Swagger *openapi3.T,

	keycloak *keycloakclient.Client,
	requiredResource string,
	requiredRole string,

	repo *messagesrepo.Repo,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(repo))
	if err != nil {
		return nil, fmt.Errorf("gethistory usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(getHistoryUseCase))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	httpErrHandler, err := errhandler.New(errhandler.NewOptions(lg, productionMode, errhandler.ResponseBuilder))
	if err != nil {
		return nil, fmt.Errorf("create errhandler: %v", err)
	}

	srv, err := serverclient.New(serverclient.NewOptions(
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
