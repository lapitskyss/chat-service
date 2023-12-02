// Code generated by MockGen. DO NOT EDIT.
// Source: service_dlq.go
//
// Generated by this command:
//
//	mockgen -source=service_dlq.go -destination=mocks/dlq_writer_mock.gen.go -package=afcverdictsprocessormocks
//
// Package afcverdictsprocessormocks is a generated GoMock package.
package afcverdictsprocessormocks

import (
	context "context"
	reflect "reflect"

	kafka "github.com/segmentio/kafka-go"
	gomock "go.uber.org/mock/gomock"
)

// MockKafkaDLQWriter is a mock of KafkaDLQWriter interface.
type MockKafkaDLQWriter struct {
	ctrl     *gomock.Controller
	recorder *MockKafkaDLQWriterMockRecorder
}

// MockKafkaDLQWriterMockRecorder is the mock recorder for MockKafkaDLQWriter.
type MockKafkaDLQWriterMockRecorder struct {
	mock *MockKafkaDLQWriter
}

// NewMockKafkaDLQWriter creates a new mock instance.
func NewMockKafkaDLQWriter(ctrl *gomock.Controller) *MockKafkaDLQWriter {
	mock := &MockKafkaDLQWriter{ctrl: ctrl}
	mock.recorder = &MockKafkaDLQWriterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockKafkaDLQWriter) EXPECT() *MockKafkaDLQWriterMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockKafkaDLQWriter) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockKafkaDLQWriterMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockKafkaDLQWriter)(nil).Close))
}

// WriteMessages mocks base method.
func (m *MockKafkaDLQWriter) WriteMessages(ctx context.Context, msgs ...kafka.Message) error {
	m.ctrl.T.Helper()
	varargs := []any{ctx}
	for _, a := range msgs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "WriteMessages", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// WriteMessages indicates an expected call of WriteMessages.
func (mr *MockKafkaDLQWriterMockRecorder) WriteMessages(ctx any, msgs ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx}, msgs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WriteMessages", reflect.TypeOf((*MockKafkaDLQWriter)(nil).WriteMessages), varargs...)
}
