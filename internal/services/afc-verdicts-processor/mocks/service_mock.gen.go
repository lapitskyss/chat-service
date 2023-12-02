// Code generated by MockGen. DO NOT EDIT.
// Source: service.go
//
// Generated by this command:
//
//	mockgen -source=service.go -destination=mocks/service_mock.gen.go -package=afcverdictsprocessormocks
//
// Package afcverdictsprocessormocks is a generated GoMock package.
package afcverdictsprocessormocks

import (
	context "context"
	reflect "reflect"
	time "time"

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

// BlockMessage mocks base method.
func (m *MockmessagesRepository) BlockMessage(ctx context.Context, msgID types.MessageID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BlockMessage", ctx, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// BlockMessage indicates an expected call of BlockMessage.
func (mr *MockmessagesRepositoryMockRecorder) BlockMessage(ctx, msgID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BlockMessage", reflect.TypeOf((*MockmessagesRepository)(nil).BlockMessage), ctx, msgID)
}

// MarkAsVisibleForManager mocks base method.
func (m *MockmessagesRepository) MarkAsVisibleForManager(ctx context.Context, msgID types.MessageID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MarkAsVisibleForManager", ctx, msgID)
	ret0, _ := ret[0].(error)
	return ret0
}

// MarkAsVisibleForManager indicates an expected call of MarkAsVisibleForManager.
func (mr *MockmessagesRepositoryMockRecorder) MarkAsVisibleForManager(ctx, msgID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MarkAsVisibleForManager", reflect.TypeOf((*MockmessagesRepository)(nil).MarkAsVisibleForManager), ctx, msgID)
}

// MockoutboxService is a mock of outboxService interface.
type MockoutboxService struct {
	ctrl     *gomock.Controller
	recorder *MockoutboxServiceMockRecorder
}

// MockoutboxServiceMockRecorder is the mock recorder for MockoutboxService.
type MockoutboxServiceMockRecorder struct {
	mock *MockoutboxService
}

// NewMockoutboxService creates a new mock instance.
func NewMockoutboxService(ctrl *gomock.Controller) *MockoutboxService {
	mock := &MockoutboxService{ctrl: ctrl}
	mock.recorder = &MockoutboxServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockoutboxService) EXPECT() *MockoutboxServiceMockRecorder {
	return m.recorder
}

// Put mocks base method.
func (m *MockoutboxService) Put(ctx context.Context, name, payload string, availableAt time.Time) (types.JobID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Put", ctx, name, payload, availableAt)
	ret0, _ := ret[0].(types.JobID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Put indicates an expected call of Put.
func (mr *MockoutboxServiceMockRecorder) Put(ctx, name, payload, availableAt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Put", reflect.TypeOf((*MockoutboxService)(nil).Put), ctx, name, payload, availableAt)
}

// Mocktransactor is a mock of transactor interface.
type Mocktransactor struct {
	ctrl     *gomock.Controller
	recorder *MocktransactorMockRecorder
}

// MocktransactorMockRecorder is the mock recorder for Mocktransactor.
type MocktransactorMockRecorder struct {
	mock *Mocktransactor
}

// NewMocktransactor creates a new mock instance.
func NewMocktransactor(ctrl *gomock.Controller) *Mocktransactor {
	mock := &Mocktransactor{ctrl: ctrl}
	mock.recorder = &MocktransactorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *Mocktransactor) EXPECT() *MocktransactorMockRecorder {
	return m.recorder
}

// RunInTx mocks base method.
func (m *Mocktransactor) RunInTx(ctx context.Context, f func(context.Context) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RunInTx", ctx, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// RunInTx indicates an expected call of RunInTx.
func (mr *MocktransactorMockRecorder) RunInTx(ctx, f any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RunInTx", reflect.TypeOf((*Mocktransactor)(nil).RunInTx), ctx, f)
}
