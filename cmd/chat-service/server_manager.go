package main

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	"github.com/lapitskyss/chat-service/internal/server"
	servermanager "github.com/lapitskyss/chat-service/internal/server-manager"
	managererrhandler "github.com/lapitskyss/chat-service/internal/server-manager/errhandler"
	managerevents "github.com/lapitskyss/chat-service/internal/server-manager/events"
	managerv1 "github.com/lapitskyss/chat-service/internal/server-manager/v1"
	"github.com/lapitskyss/chat-service/internal/server/errhandler"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	managerload "github.com/lapitskyss/chat-service/internal/services/manager-load"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/store"
	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	closechat "github.com/lapitskyss/chat-service/internal/usecases/manager/close-chat"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/send-message"
	managertypingmessage "github.com/lapitskyss/chat-service/internal/usecases/manager/typing-message"
	websocketstream "github.com/lapitskyss/chat-service/internal/websocket-stream"
	managerhandler "github.com/lapitskyss/chat-service/internal/websocket-stream/manager-handler"
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

	db *store.Database,
	chatRepo *chatsrepo.Repo,
	msgRepo *messagesrepo.Repo,
	problemRepo *problemsrepo.Repo,

	eventStream eventstream.EventStream,
	managerLoadSvc *managerload.Service,
	managerPool managerpool.Pool,
	outboxSvc *outbox.Service,
) (*server.Server, error) {
	lg := zap.L().Named(nameServerManager)

	canReceiveProblemUseCase, err := canreceiveproblems.New(canreceiveproblems.NewOptions(
		managerLoadSvc,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("canreceiveproblems usecase: %v", err)
	}
	closeChatUseCase, err := closechat.New(closechat.NewOptions(
		msgRepo,
		outboxSvc,
		problemRepo,
		db,
	))
	if err != nil {
		return nil, fmt.Errorf("closechat usecase: %v", err)
	}
	freeHandsUseCase, err := freehands.New(freehands.NewOptions(
		managerLoadSvc,
		managerPool,
	))
	if err != nil {
		return nil, fmt.Errorf("freehands usecase: %v", err)
	}
	getChatHistoryUseCase, err := getchathistory.New(getchathistory.NewOptions(
		msgRepo,
		problemRepo,
	))
	if err != nil {
		return nil, fmt.Errorf("getchats usecase: %v", err)
	}
	getChatsUseCase, err := getchats.New(getchats.NewOptions(
		chatRepo,
	))
	if err != nil {
		return nil, fmt.Errorf("getchats usecase: %v", err)
	}
	sendMessageUseCase, err := sendmessage.New(sendmessage.NewOptions(
		msgRepo,
		outboxSvc,
		problemRepo,
		db,
	))
	if err != nil {
		return nil, fmt.Errorf("sendmessage usecase: %v", err)
	}
	typingMessageUseCase, err := managertypingmessage.New(managertypingmessage.NewOptions(
		chatRepo,
		problemRepo,
		eventStream,
	))
	if err != nil {
		return nil, fmt.Errorf("managertypingmessage usecase: %v", err)
	}

	v1Handlers, err := managerv1.NewHandlers(managerv1.NewOptions(
		canReceiveProblemUseCase,
		closeChatUseCase,
		freeHandsUseCase,
		getChatHistoryUseCase,
		getChatsUseCase,
		sendMessageUseCase,
	))
	if err != nil {
		return nil, fmt.Errorf("create v1 manager handlers: %v", err)
	}

	shutdownCh := make(chan struct{})
	shutdown := func() {
		close(shutdownCh)
	}

	wsReadHandler, err := managerhandler.New(managerhandler.NewOptions(
		managerhandler.JSONEventReader{},
		typingMessageUseCase,
	))
	if err != nil {
		return nil, fmt.Errorf("create ws manager read handler: %v", err)
	}

	wsUpgrader := websocketstream.NewUpgrader(allowOrigins, secWsProtocol)
	wsHandler, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		lg,
		eventStream,
		managerevents.Adapter{},
		websocketstream.JSONEventWriter{},
		wsUpgrader,
		wsReadHandler,
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
		shutdown,
	))
	if err != nil {
		return nil, fmt.Errorf("build server: %v", err)
	}

	return srv, nil
}
