// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=gethistorymocks
//
// Package gethistorymocks is a generated GoMock package.
package gethistorymocks

import (
	context "context"
	reflect "reflect"

	messagesrepo "github.com/lapitskyss/chat-service/internal/repositories/messages"
	types "github.com/lapitskyss/chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

// MockmessagesRepository is a mock of messagesRepository interface.
type MockmessagesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockmessagesRepositoryMockRecorder
}

// MockmessagesRepositoryMockRecorder is the mock recorder for MockmessagesRepository.
type MockmessagesRepositoryMockRecorder struct {
	mock *MockmessagesRepository
}

// NewMockmessagesRepository creates a new mock instance.
func NewMockmessagesRepository(ctrl *gomock.Controller) *MockmessagesRepository {
	mock := &MockmessagesRepository{ctrl: ctrl}
	mock.recorder = &MockmessagesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmessagesRepository) EXPECT() *MockmessagesRepositoryMockRecorder {
	return m.recorder
}

// GetClientChatMessages mocks base method.
func (m *MockmessagesRepository) GetClientChatMessages(ctx context.Context, clientID types.UserID, pageSize int, cursor *messagesrepo.Cursor) ([]messagesrepo.Message, *messagesrepo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientChatMessages", ctx, clientID, pageSize, cursor)
	ret0, _ := ret[0].([]messagesrepo.Message)
	ret1, _ := ret[1].(*messagesrepo.Cursor)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetClientChatMessages indicates an expected call of GetClientChatMessages.
func (mr *MockmessagesRepositoryMockRecorder) GetClientChatMessages(ctx, clientID, pageSize, cursor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientChatMessages", reflect.TypeOf((*MockmessagesRepository)(nil).GetClientChatMessages), ctx, clientID, pageSize, cursor)
}
