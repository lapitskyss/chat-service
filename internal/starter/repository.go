package starter

import (
	"fmt"

	"github.com/google/wire"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	jobsrepo "github.com/lapitskyss/chat-service/internal/repositories/jobs"
	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	"github.com/lapitskyss/chat-service/internal/store"
)

//nolint:unused
var repositorySet = wire.NewSet(
	provideChatRepo,
	provideMsgRepo,
	provideProblemRepo,
	provideJobsRepo,
)

func provideChatRepo(db *store.Database) (*chatsrepo.Repo, error) {
	chatRepo, err := chatsrepo.New(chatsrepo.NewOptions(db))
	if err != nil {
		return nil, fmt.Errorf("chats repository: %v", err)
	}
	return chatRepo, nil
}

func provideMsgRepo(db *store.Database) (*messagesrepo.Repo, error) {
	msgRepo, err := messagesrepo.New(messagesrepo.NewOptions(db))
	if err != nil {
		return nil, fmt.Errorf("messages repository: %v", err)
	}
	return msgRepo, nil
}

func provideProblemRepo(db *store.Database) (*problemsrepo.Repo, error) {
	problemRepo, err := problemsrepo.New(problemsrepo.NewOptions(db))
	if err != nil {
		return nil, fmt.Errorf("messages repository: %v", err)
	}
	return problemRepo, nil
}

func provideJobsRepo(db *store.Database) (*jobsrepo.Repo, error) {
	jobsRepo, err := jobsrepo.New(jobsrepo.NewOptions(db))
	if err != nil {
		return nil, fmt.Errorf("jobs repository: %v", err)
	}
	return jobsRepo, nil
}
