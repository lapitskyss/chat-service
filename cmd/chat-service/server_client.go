package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	serverclient "github.com/lapitskyss/chat-service/internal/server-client"
	"github.com/lapitskyss/chat-service/internal/server-client/errhandler"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/store"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
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

	db *store.Database,
	chatRepo *chatsrepo.Repo,
	msgRepo *messagesrepo.Repo,
	problemRepo *problemsrepo.Repo,

	outboxSvc *outbox.Service,
) (*serverclient.Server, error) {
	lg := zap.L().Named(nameServerClient)

	getHistoryUseCase, err := gethistory.New(gethistory.NewOptions(msgRepo))
	if err != nil {
		return nil, fmt.Errorf("gethistory usecase: %v", err)
	}

	sendMessageUseCase, err := sendmessage.New(sendmessage.NewOptions(
		chatRepo,
		msgRepo,
		outboxSvc,
		problemRepo,
		db,
	))
	if err != nil {
		return nil, fmt.Errorf("gethistory usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(
		getHistoryUseCase,
		sendMessageUseCase,
	))
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
