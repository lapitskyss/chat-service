package starter

import (
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"

	"github.com/lapitskyss/chat-service/internal/config"
	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	afcverdictsprocessor "github.com/lapitskyss/chat-service/internal/services/afc-verdicts-processor"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	inmemeventstream "github.com/lapitskyss/chat-service/internal/services/event-stream/in-mem"
	managerload "github.com/lapitskyss/chat-service/internal/services/manager-load"
	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	inmemmanagerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool/in-mem"
	managerscheduler "github.com/lapitskyss/chat-service/internal/services/manager-scheduler"
	msgproducer "github.com/lapitskyss/chat-service/internal/services/msg-producer"
	"github.com/lapitskyss/chat-service/internal/services/outbox"
	clientmessageblockedjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/client-message-blocked"
	clientmessagesentjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/client-message-sent"
	managerassignedtoproblemjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/manager-assigned-to-problem"
	managerclosechatjob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/manager-close-chat"
	sendclientmessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-client-message"
	sendmanagermessagejob "github.com/lapitskyss/chat-service/internal/services/outbox/jobs/send-manager-message"
	"github.com/lapitskyss/chat-service/internal/store"
)

//nolint:unused
var servicesSet = wire.NewSet(
	provideMsgProducer,
	provideInMemoryEventStream,
	provideInMemoryManagerPool,
	provideOutbox,
	provideManagerLoad,
	provideManagerScheduler,
	provideAFCVerdictsProcessor,
)

func provideMsgProducer(cfg config.Config) (*msgproducer.Service, error) {
	msgProducer, err := msgproducer.New(msgproducer.NewOptions(
		msgproducer.NewKafkaWriter(
			cfg.Services.MsgProducer.Brokers,
			cfg.Services.MsgProducer.Topic,
			cfg.Services.MsgProducer.BatchSize,
		),
		msgproducer.WithEncryptKey(cfg.Services.MsgProducer.EncryptKey),
	))
	if err != nil {
		return nil, fmt.Errorf("message producer service: %v", err)
	}
	return msgProducer, nil
}

func provideInMemoryEventStream(log *zap.Logger) (eventstream.EventStream, func()) {
	eventStream := inmemeventstream.New()

	cleanup := func() {
		err := eventStream.Close()
		if err != nil {
			log.Error("close event stream", zap.Error(err))
		}
	}

	return eventStream, cleanup
}

func provideInMemoryManagerPool() managerpool.Pool {
	return inmemmanagerpool.New()
}

func provideManagerLoad(cfg config.Config, problemRepo *problemsrepo.Repo) (*managerload.Service, error) {
	managerLoad, err := managerload.New(managerload.NewOptions(
		cfg.Services.ManagerLoad.MaxProblemsAtTime,
		problemRepo,
	))
	if err != nil {
		return nil, fmt.Errorf("manager load service: %v", err)
	}
	return managerLoad, nil
}

func provideOutbox(
	cfg config.Config,
	chatRepo *chatsrepo.Repo,
	msgRepo *messagesrepo.Repo,
	jobsRepo *jobsrepo.Repo,
	db *store.Database,

	msgProducer *msgproducer.Service,
	eventsStream eventstream.EventStream,
	managerLoad *managerload.Service,
) (*outbox.Service, error) {
	outBox, err := outbox.New(outbox.NewOptions(
		cfg.Services.Outbox.Workers,
		cfg.Services.Outbox.IdleTime,
		cfg.Services.Outbox.ReserveFor,
		jobsRepo,
		db,
	))
	if err != nil {
		return nil, fmt.Errorf("outbox service: %v", err)
	}

	sendClientMessageJob, err := sendclientmessagejob.New(sendclientmessagejob.NewOptions(
		msgProducer,
		msgRepo,
		eventsStream,
	))
	if err != nil {
		return nil, fmt.Errorf("send client message job: %v", err)
	}
	err = outBox.RegisterJob(sendClientMessageJob)
	if err != nil {
		return nil, fmt.Errorf("register send client message job: %v", err)
	}

	clientMessageBlockedJob, err := clientmessageblockedjob.New(clientmessageblockedjob.NewOptions(
		msgRepo,
		eventsStream,
	))
	if err != nil {
		return nil, fmt.Errorf("client message blocked job: %v", err)
	}
	err = outBox.RegisterJob(clientMessageBlockedJob)
	if err != nil {
		return nil, fmt.Errorf("register client message blocked job: %v", err)
	}

	clientMessageSentJob, err := clientmessagesentjob.New(clientmessagesentjob.NewOptions(
		msgRepo,
		eventsStream,
	))
	if err != nil {
		return nil, fmt.Errorf("client message sent job: %v", err)
	}
	err = outBox.RegisterJob(clientMessageSentJob)
	if err != nil {
		return nil, fmt.Errorf("register client message sent job: %v", err)
	}

	managerAssignedToProblemJob, err := managerassignedtoproblemjob.New(managerassignedtoproblemjob.NewOptions(
		managerLoad,
		msgRepo,
		eventsStream,
	))
	if err != nil {
		return nil, fmt.Errorf("manager assigned to problem job: %v", err)
	}
	err = outBox.RegisterJob(managerAssignedToProblemJob)
	if err != nil {
		return nil, fmt.Errorf("register manager assigned to problem job: %v", err)
	}

	sendManagerMessageJob, err := sendmanagermessagejob.New(sendmanagermessagejob.NewOptions(
		msgRepo,
		chatRepo,
		eventsStream,
		msgProducer,
	))
	if err != nil {
		return nil, fmt.Errorf("send manager message job: %v", err)
	}
	err = outBox.RegisterJob(sendManagerMessageJob)
	if err != nil {
		return nil, fmt.Errorf("register send manager message job: %v", err)
	}
	managerCloseChatJob, err := managerclosechatjob.New(managerclosechatjob.NewOptions(
		msgRepo,
		eventsStream,
		managerLoad,
	))
	if err != nil {
		return nil, fmt.Errorf("manager close chat job: %v", err)
	}
	err = outBox.RegisterJob(managerCloseChatJob)
	if err != nil {
		return nil, fmt.Errorf("manager close chat job: %v", err)
	}

	return outBox, nil
}

func provideManagerScheduler(
	cfg config.Config,
	managerPool managerpool.Pool,
	msgRepo *messagesrepo.Repo,
	outBox *outbox.Service,
	problemRepo *problemsrepo.Repo,
	db *store.Database,
) (*managerscheduler.Service, error) {
	managerScheduler, err := managerscheduler.New(managerscheduler.NewOptions(
		cfg.Services.ManagerScheduler.Period,
		managerPool,
		msgRepo,
		outBox,
		problemRepo,
		db,
	))
	if err != nil {
		return nil, fmt.Errorf("manager scheduler service: %v", err)
	}
	return managerScheduler, nil
}

func provideAFCVerdictsProcessor(
	cfg config.Config,
	db *store.Database,
	msgRepo *messagesrepo.Repo,
	outBox *outbox.Service,
) (*afcverdictsprocessor.Service, error) {
	afcVerdictProcessor, err := afcverdictsprocessor.New(afcverdictsprocessor.NewOptions(
		cfg.Services.AFCVerdictsProcessor.Brokers,
		cfg.Services.AFCVerdictsProcessor.Consumers,
		cfg.Services.AFCVerdictsProcessor.ConsumerGroup,
		cfg.Services.AFCVerdictsProcessor.VerdictsTopic,
		afcverdictsprocessor.NewKafkaReader,
		afcverdictsprocessor.NewKafkaDLQWriter(
			cfg.Services.AFCVerdictsProcessor.Brokers,
			cfg.Services.AFCVerdictsProcessor.VerdictsDlqTopic,
		),
		db,
		msgRepo,
		outBox,
		afcverdictsprocessor.WithVerdictsSignKey(cfg.Services.AFCVerdictsProcessor.VerdictsSigningPublicKey),
		afcverdictsprocessor.WithProcessBatchSize(cfg.Services.AFCVerdictsProcessor.BatchSize),
	))
	if err != nil {
		return nil, fmt.Errorf("AFC verdict processor: %v", err)
	}
	return afcVerdictProcessor, nil
}
