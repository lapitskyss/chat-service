// Code generated by MockGen. DO NOT EDIT.
// Source: usecase.go
//
// Generated by this command:
//
//	mockgen -source=usecase.go -destination=mocks/usecase_mock.gen.go -package=getchatsmocks
//
// Package getchatsmocks is a generated GoMock package.
package getchatsmocks

import (
	context "context"
	reflect "reflect"

	chatsrepo "github.com/lapitskyss/chat-service/internal/repositories/chats"
	types "github.com/lapitskyss/chat-service/internal/types"
	gomock "go.uber.org/mock/gomock"
)

// MockchatsRepository is a mock of chatsRepository interface.
type MockchatsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockchatsRepositoryMockRecorder
}

// MockchatsRepositoryMockRecorder is the mock recorder for MockchatsRepository.
type MockchatsRepositoryMockRecorder struct {
	mock *MockchatsRepository
}

// NewMockchatsRepository creates a new mock instance.
func NewMockchatsRepository(ctrl *gomock.Controller) *MockchatsRepository {
	mock := &MockchatsRepository{ctrl: ctrl}
	mock.recorder = &MockchatsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockchatsRepository) EXPECT() *MockchatsRepositoryMockRecorder {
	return m.recorder
}

// AllWithOpenProblemsForManager mocks base method.
func (m *MockchatsRepository) AllWithOpenProblemsForManager(ctx context.Context, managerID types.UserID) ([]chatsrepo.Chat, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AllWithOpenProblemsForManager", ctx, managerID)
	ret0, _ := ret[0].([]chatsrepo.Chat)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AllWithOpenProblemsForManager indicates an expected call of AllWithOpenProblemsForManager.
func (mr *MockchatsRepositoryMockRecorder) AllWithOpenProblemsForManager(ctx, managerID any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AllWithOpenProblemsForManager", reflect.TypeOf((*MockchatsRepository)(nil).AllWithOpenProblemsForManager), ctx, managerID)
}
