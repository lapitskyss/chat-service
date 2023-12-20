// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=closechatmocks
//
// Package closechatmocks is a generated GoMock package.
package closechatmocks

import (
	context "context"
	reflect "reflect"
	time "time"

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

// CreateServiceMsg mocks base method.
func (m *MockmessagesRepository) CreateServiceMsg(ctx context.Context, reqID types.RequestID, problemID types.ProblemID, chatID types.ChatID, msgBody string, visibleForClient, visibleForManager bool) (*messagesrepo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateServiceMsg", ctx, reqID, problemID, chatID, msgBody, visibleForClient, visibleForManager)
	ret0, _ := ret[0].(*messagesrepo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateServiceMsg indicates an expected call of CreateServiceMsg.
func (mr *MockmessagesRepositoryMockRecorder) CreateServiceMsg(ctx, reqID, problemID, chatID, msgBody, visibleForClient, visibleForManager any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateServiceMsg", reflect.TypeOf((*MockmessagesRepository)(nil).CreateServiceMsg), ctx, reqID, problemID, chatID, msgBody, visibleForClient, visibleForManager)
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

// ResolveProblem mocks base method.
func (m *MockproblemsRepository) ResolveProblem(ctx context.Context, problemID types.ProblemID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResolveProblem", ctx, problemID)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResolveProblem indicates an expected call of ResolveProblem.
func (mr *MockproblemsRepositoryMockRecorder) ResolveProblem(ctx, problemID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResolveProblem", reflect.TypeOf((*MockproblemsRepository)(nil).ResolveProblem), ctx, problemID)
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