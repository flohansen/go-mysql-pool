// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/flohansen/go-mysql-pool (interfaces: Conn)
//
// Generated by this command:
//
//	mockgen -destination=mocks/conn.go -package=mocks github.com/flohansen/go-mysql-pool Conn
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockConn is a mock of Conn interface.
type MockConn struct {
	ctrl     *gomock.Controller
	recorder *MockConnMockRecorder
	isgomock struct{}
}

// MockConnMockRecorder is the mock recorder for MockConn.
type MockConnMockRecorder struct {
	mock *MockConn
}

// NewMockConn creates a new mock instance.
func NewMockConn(ctrl *gomock.Controller) *MockConn {
	mock := &MockConn{ctrl: ctrl}
	mock.recorder = &MockConnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConn) EXPECT() *MockConnMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockConn) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockConnMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockConn)(nil).Close))
}

// ExecContext mocks base method.
func (m *MockConn) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockConnMockRecorder) ExecContext(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockConn)(nil).ExecContext), varargs...)
}

// QueryContext mocks base method.
func (m *MockConn) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	m.ctrl.T.Helper()
	varargs := []any{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "QueryContext", varargs...)
	ret0, _ := ret[0].(*sql.Rows)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// QueryContext indicates an expected call of QueryContext.
func (mr *MockConnMockRecorder) QueryContext(ctx, query any, args ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueryContext", reflect.TypeOf((*MockConn)(nil).QueryContext), varargs...)
}
