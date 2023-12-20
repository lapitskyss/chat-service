// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=getchathistorymocks
//
// Package getchathistorymocks is a generated GoMock package.
package getchathistorymocks

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

// GetProblemMessages mocks base method.
func (m *MockmessagesRepository) GetProblemMessages(ctx context.Context, problemID types.ProblemID, pageSize int, cursor *messagesrepo.Cursor) ([]messagesrepo.Message, *messagesrepo.Cursor, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProblemMessages", ctx, problemID, pageSize, cursor)
	ret0, _ := ret[0].([]messagesrepo.Message)
	ret1, _ := ret[1].(*messagesrepo.Cursor)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetProblemMessages indicates an expected call of GetProblemMessages.
func (mr *MockmessagesRepositoryMockRecorder) GetProblemMessages(ctx, problemID, pageSize, cursor any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProblemMessages", reflect.TypeOf((*MockmessagesRepository)(nil).GetProblemMessages), ctx, problemID, pageSize, cursor)
}

// MockproblemsRepository is a mock of problemsRepository interface.
type MockproblemsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockproblemsRepositoryMockRecorder
}

// MockproblemsRepositoryMockRecorder is the mock recorder for MockproblemsRepository.
type MockproblemsRepositoryMockRecorder struct {
	mock *MockproblemsRepository
}

// NewMockproblemsRepository creates a new mock instance.
func NewMockproblemsRepository(ctrl *gomock.Controller) *MockproblemsRepository {
	mock := &MockproblemsRepository{ctrl: ctrl}
	mock.recorder = &MockproblemsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockproblemsRepository) EXPECT() *MockproblemsRepositoryMockRecorder {
	return m.recorder
}

// GetAssignedProblemID mocks base method.
func (m *MockproblemsRepository) GetAssignedProblemID(ctx context.Context, managerID types.UserID, chatID types.ChatID) (types.ProblemID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAssignedProblemID", ctx, managerID, chatID)
	ret0, _ := ret[0].(types.ProblemID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAssignedProblemID indicates an expected call of GetAssignedProblemID.
func (mr *MockproblemsRepositoryMockRecorder) GetAssignedProblemID(ctx, managerID, chatID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAssignedProblemID", reflect.TypeOf((*MockproblemsRepository)(nil).GetAssignedProblemID), ctx, managerID, chatID)
}