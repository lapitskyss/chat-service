package starter

import (
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"

	keycloakclient "github.com/lapitskyss/chat-service/internal/clients/keycloak"
	"github.com/lapitskyss/chat-service/internal/config"
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	"github.com/lapitskyss/chat-service/internal/server"
	serverclient "github.com/lapitskyss/chat-service/internal/server-client"
	clienterrhandler "github.com/lapitskyss/chat-service/internal/server-client/errhandler"
	clientevents "github.com/lapitskyss/chat-service/internal/server-client/events"
	clientv1 "github.com/lapitskyss/chat-service/internal/server-client/v1"
	"github.com/lapitskyss/chat-service/internal/server/errhandler"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	"github.com/lapitskyss/chat-service/internal/store"
	gethistory "github.com/lapitskyss/chat-service/internal/usecases/client/get-history"
	sendmessage "github.com/lapitskyss/chat-service/internal/usecases/client/send-message"
	clienttypingmessage "github.com/lapitskyss/chat-service/internal/usecases/client/typing-message"
	websocketstream "github.com/lapitskyss/chat-service/internal/websocket-stream"
	clienthandler "github.com/lapitskyss/chat-service/internal/websocket-stream/client-handler"
)

type ServerClient *server.Server

//nolint:unused
var serverClientSet = wire.NewSet(
	provideServerClient,
)

const nameServerClient = "server-client"

func provideServerClient(
	cfg config.Config,

	db *store.Database,
	chatRepo *chatsrepo.Repo,
	msgRepo *messagesrepo.Repo,
	problemRepo *problemsrepo.Repo,

	eventStream eventstream.EventStream,
	outboxSvc *outbox.Service,

	keycloak *keycloakclient.Client,
	v1Swagger ClientV1Swagger,
) (ServerClient, error) {
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
	typingMessageUseCase, err := clienttypingmessage.New(clienttypingmessage.NewOptions(
		problemRepo,
		eventStream,
	))
	if err != nil {
		return nil, fmt.Errorf("clienttypingmessage usecase: %v", err)
	}

	v1Handlers, err := clientv1.NewHandlers(clientv1.NewOptions(
		getHistoryUseCase,
		sendMessageUseCase,
	))
	if err != nil {
		return nil, fmt.Errorf("create v1 handlers: %v", err)
	}

	wsReadHandler, err := clienthandler.New(clienthandler.NewOptions(
		clienthandler.JSONEventReader{},
		typingMessageUseCase,
	))
	if err != nil {
		return nil, fmt.Errorf("create ws client read handler: %v", err)
	}

	shutdownCh := make(chan struct{})
	shutdown := func() {
		close(shutdownCh)
	}

	wsUpgrader := websocketstream.NewUpgrader(cfg.Servers.Client.AllowOrigins, cfg.Servers.Client.SecWsProtocol)
	wsHandler, err := websocketstream.NewHTTPHandler(websocketstream.NewOptions(
		lg,
		eventStream,
		clientevents.Adapter{},
		websocketstream.JSONEventWriter{},
		wsUpgrader,
		wsReadHandler,
		shutdownCh,
	))
	if err != nil {
		return nil, fmt.Errorf("websock etstream handler: %v", err)
	}

	httpErrHandler, err := errhandler.New(
		errhandler.NewOptions(lg, cfg.Global.IsProd(), clienterrhandler.ResponseBuilder),
	)
	if err != nil {
		return nil, fmt.Errorf("create errhandler: %v", err)
	}

	srv, err := serverclient.New(serverclient.NewOptions(
		lg,
		cfg.Servers.Client.Addr,
		cfg.Servers.Client.AllowOrigins,
		keycloak,
		cfg.Servers.Client.RequiredAccess.Resource,
		cfg.Servers.Client.RequiredAccess.Role,
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
