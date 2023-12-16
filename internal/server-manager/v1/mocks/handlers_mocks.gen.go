// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go
//
// Generated by this command:
//
//	mockgen -source=handlers.go -destination=mocks/handlers_mocks.gen.go -package=managerv1mocks
//
// Package managerv1mocks is a generated GoMock package.
package managerv1mocks

import (
	context "context"
	reflect "reflect"

	canreceiveproblems "github.com/lapitskyss/chat-service/internal/usecases/manager/can-receive-problems"
	freehands "github.com/lapitskyss/chat-service/internal/usecases/manager/free-hands"
	getchathistory "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chat-history"
	getchats "github.com/lapitskyss/chat-service/internal/usecases/manager/get-chats"
	gomock "go.uber.org/mock/gomock"
)

// MockcanReceiveProblemsUseCase is a mock of canReceiveProblemsUseCase interface.
type MockcanReceiveProblemsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockcanReceiveProblemsUseCaseMockRecorder
}

// MockcanReceiveProblemsUseCaseMockRecorder is the mock recorder for MockcanReceiveProblemsUseCase.
type MockcanReceiveProblemsUseCaseMockRecorder struct {
	mock *MockcanReceiveProblemsUseCase
}

// NewMockcanReceiveProblemsUseCase creates a new mock instance.
func NewMockcanReceiveProblemsUseCase(ctrl *gomock.Controller) *MockcanReceiveProblemsUseCase {
	mock := &MockcanReceiveProblemsUseCase{ctrl: ctrl}
	mock.recorder = &MockcanReceiveProblemsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcanReceiveProblemsUseCase) EXPECT() *MockcanReceiveProblemsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockcanReceiveProblemsUseCase) Handle(ctx context.Context, req canreceiveproblems.Request) (canreceiveproblems.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(canreceiveproblems.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockcanReceiveProblemsUseCaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockcanReceiveProblemsUseCase)(nil).Handle), ctx, req)
}

// MockfreeHandsUseCase is a mock of freeHandsUseCase interface.
type MockfreeHandsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockfreeHandsUseCaseMockRecorder
}

// MockfreeHandsUseCaseMockRecorder is the mock recorder for MockfreeHandsUseCase.
type MockfreeHandsUseCaseMockRecorder struct {
	mock *MockfreeHandsUseCase
}

// NewMockfreeHandsUseCase creates a new mock instance.
func NewMockfreeHandsUseCase(ctrl *gomock.Controller) *MockfreeHandsUseCase {
	mock := &MockfreeHandsUseCase{ctrl: ctrl}
	mock.recorder = &MockfreeHandsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockfreeHandsUseCase) EXPECT() *MockfreeHandsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockfreeHandsUseCase) Handle(ctx context.Context, req freehands.Request) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockfreeHandsUseCaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockfreeHandsUseCase)(nil).Handle), ctx, req)
}

// MockgetChatHistoryUseCase is a mock of getChatHistoryUseCase interface.
type MockgetChatHistoryUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetChatHistoryUseCaseMockRecorder
}

// MockgetChatHistoryUseCaseMockRecorder is the mock recorder for MockgetChatHistoryUseCase.
type MockgetChatHistoryUseCaseMockRecorder struct {
	mock *MockgetChatHistoryUseCase
}

// NewMockgetChatHistoryUseCase creates a new mock instance.
func NewMockgetChatHistoryUseCase(ctrl *gomock.Controller) *MockgetChatHistoryUseCase {
	mock := &MockgetChatHistoryUseCase{ctrl: ctrl}
	mock.recorder = &MockgetChatHistoryUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetChatHistoryUseCase) EXPECT() *MockgetChatHistoryUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetChatHistoryUseCase) Handle(ctx context.Context, req getchathistory.Request) (getchathistory.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getchathistory.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetChatHistoryUseCaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetChatHistoryUseCase)(nil).Handle), ctx, req)
}

// MockgetChatsUseCase is a mock of getChatsUseCase interface.
type MockgetChatsUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockgetChatsUseCaseMockRecorder
}

// MockgetChatsUseCaseMockRecorder is the mock recorder for MockgetChatsUseCase.
type MockgetChatsUseCaseMockRecorder struct {
	mock *MockgetChatsUseCase
}

// NewMockgetChatsUseCase creates a new mock instance.
func NewMockgetChatsUseCase(ctrl *gomock.Controller) *MockgetChatsUseCase {
	mock := &MockgetChatsUseCase{ctrl: ctrl}
	mock.recorder = &MockgetChatsUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockgetChatsUseCase) EXPECT() *MockgetChatsUseCaseMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockgetChatsUseCase) Handle(ctx context.Context, req getchats.Request) (getchats.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", ctx, req)
	ret0, _ := ret[0].(getchats.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Handle indicates an expected call of Handle.
func (mr *MockgetChatsUseCaseMockRecorder) Handle(ctx, req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockgetChatsUseCase)(nil).Handle), ctx, req)
}
