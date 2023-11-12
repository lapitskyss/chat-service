package inmemmanagerpool

import (
	"context"
	"sync"

	managerpool "github.com/lapitskyss/chat-service/internal/services/manager-pool"
	"github.com/lapitskyss/chat-service/internal/types"
)

type Service struct {
	head, tail *node
	items      map[types.UserID]struct{}

	mu sync.Mutex
}

type node struct {
	next, prev *node
	managerID  types.UserID
}

func (s *Service) Get(_ context.Context) (types.UserID, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.tail == nil {
		return types.UserIDNil, managerpool.ErrNoAvailableManagers
	}

	result := s.tail.managerID
	delete(s.items, s.tail.managerID)

	if s.tail.prev != nil {
		s.tail = s.tail.prev
		s.tail.next = nil
	} else {
		s.tail = nil
		s.head = nil
	}

	return result, nil
}

func (s *Service) Put(_ context.Context, managerID types.UserID) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exist := s.items[managerID]; exist {
		return nil
	}

	n := &node{
		managerID: managerID,
		next:      s.head,
	}

	if s.head != nil {
		s.head.prev = n
	}
	if s.tail == nil {
		s.tail = n
	}

	s.head = n
	s.items[managerID] = struct{}{}

	return nil
}

func (s *Service) Contains(_ context.Context, managerID types.UserID) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.items[managerID]
	return exist, nil
}

func (s *Service) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return len(s.items)
}

func (s *Service) Close() error {
	return nil
}

func New() *Service {
	return &Service{
		head:  nil,
		tail:  nil,
		items: make(map[types.UserID]struct{}),
	}
}
