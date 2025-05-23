// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/moevm/nosql1h25-writer/backend/internal/service/auth (interfaces: Service)
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

	uuid "github.com/google/uuid"
	entity "github.com/moevm/nosql1h25-writer/backend/internal/entity"
	auth "github.com/moevm/nosql1h25-writer/backend/internal/service/auth"
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

// Login mocks base method.
func (m *MockService) Login(ctx context.Context, email, password string) (entity.AuthData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(entity.AuthData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockServiceMockRecorder) Login(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockService)(nil).Login), ctx, email, password)
}

// Logout mocks base method.
func (m *MockService) Logout(ctx context.Context, refreshToken uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", ctx, refreshToken)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout.
func (mr *MockServiceMockRecorder) Logout(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockService)(nil).Logout), ctx, refreshToken)
}

// ParseToken mocks base method.
func (m *MockService) ParseToken(tokenString string) (*entity.AccessTokenClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", tokenString)
	ret0, _ := ret[0].(*entity.AccessTokenClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockServiceMockRecorder) ParseToken(tokenString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockService)(nil).ParseToken), tokenString)
}

// Refresh mocks base method.
func (m *MockService) Refresh(ctx context.Context, refreshToken uuid.UUID) (entity.AuthData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", ctx, refreshToken)
	ret0, _ := ret[0].(entity.AuthData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Refresh indicates an expected call of Refresh.
func (mr *MockServiceMockRecorder) Refresh(ctx, refreshToken any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockService)(nil).Refresh), ctx, refreshToken)
}

// Register mocks base method.
func (m *MockService) Register(ctx context.Context, in auth.RegisterIn) (entity.AuthData, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, in)
	ret0, _ := ret[0].(entity.AuthData)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockServiceMockRecorder) Register(ctx, in any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockService)(nil).Register), ctx, in)
}
