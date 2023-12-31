// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=clienttypingmessagemocks
//
// Package clienttypingmessagemocks is a generated GoMock package.
package clienttypingmessagemocks

import (
	context "context"
	reflect "reflect"

	problemsrepo "github.com/lapitskyss/chat-service/internal/repositories/problems"
	eventstream "github.com/lapitskyss/chat-service/internal/services/event-stream"
	types "github.com/lapitskyss/chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

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

// GetClientOpenProblem mocks base method.
func (m *MockproblemsRepository) GetClientOpenProblem(ctx context.Context, clientID types.UserID) (*problemsrepo.Problem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientOpenProblem", ctx, clientID)
	ret0, _ := ret[0].(*problemsrepo.Problem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClientOpenProblem indicates an expected call of GetClientOpenProblem.
func (mr *MockproblemsRepositoryMockRecorder) GetClientOpenProblem(ctx, clientID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientOpenProblem", reflect.TypeOf((*MockproblemsRepository)(nil).GetClientOpenProblem), ctx, clientID)
}

// MockeventStream is a mock of eventStream interface.
type MockeventStream struct {
	ctrl     *gomock.Controller
	recorder *MockeventStreamMockRecorder
}

// MockeventStreamMockRecorder is the mock recorder for MockeventStream.
type MockeventStreamMockRecorder struct {
	mock *MockeventStream
}

// NewMockeventStream creates a new mock instance.
func NewMockeventStream(ctrl *gomock.Controller) *MockeventStream {
	mock := &MockeventStream{ctrl: ctrl}
	mock.recorder = &MockeventStreamMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockeventStream) EXPECT() *MockeventStreamMockRecorder {
	return m.recorder
}

// Publish mocks base method.
func (m *MockeventStream) Publish(ctx context.Context, userID types.UserID, event eventstream.Event) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, userID, event)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockeventStreamMockRecorder) Publish(ctx, userID, event any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockeventStream)(nil).Publish), ctx, userID, event)
}
