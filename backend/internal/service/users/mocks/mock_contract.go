// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/moevm/nosql1h25-writer/backend/internal/service/users (interfaces: Service)
//
// Generated by this command:
//
//	mockgen -destination mocks/mock_contract.go -package=mocks . Service
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	users "github.com/moevm/nosql1h25-writer/backend/internal/service/users"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
	isgomock struct{}
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// UpdateBalance mocks base method.
func (m *MockService) UpdateBalance(ctx context.Context, userID primitive.ObjectID, op users.OperationType, amount int) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", ctx, userID, op, amount)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockServiceMockRecorder) UpdateBalance(ctx, userID, op, amount any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockService)(nil).UpdateBalance), ctx, userID, op, amount)
}
