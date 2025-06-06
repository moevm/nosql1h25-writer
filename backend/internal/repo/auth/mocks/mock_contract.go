// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/moevm/nosql1h25-writer/backend/internal/repo/auth (interfaces: Repo)
//
// Generated by this command:
//
//	mockgen -destination mocks/mock_contract.go -package=mocks . Repo
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	uuid "github.com/google/uuid"
	entity "github.com/moevm/nosql1h25-writer/backend/internal/entity"
	primitive "go.mongodb.org/mongo-driver/bson/primitive"
	gomock "go.uber.org/mock/gomock"
)

// MockRepo is a mock of Repo interface.
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
	isgomock struct{}
}

// MockRepoMockRecorder is the mock recorder for MockRepo.
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance.
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// CreateSession mocks base method.
func (m *MockRepo) CreateSession(ctx context.Context, userID primitive.ObjectID, ttl time.Duration) (entity.RefreshSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateSession", ctx, userID, ttl)
	ret0, _ := ret[0].(entity.RefreshSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateSession indicates an expected call of CreateSession.
func (mr *MockRepoMockRecorder) CreateSession(ctx, userID, ttl any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateSession", reflect.TypeOf((*MockRepo)(nil).CreateSession), ctx, userID, ttl)
}

// DeleteByToken mocks base method.
func (m *MockRepo) DeleteByToken(ctx context.Context, token uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByToken", ctx, token)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByToken indicates an expected call of DeleteByToken.
func (mr *MockRepoMockRecorder) DeleteByToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByToken", reflect.TypeOf((*MockRepo)(nil).DeleteByToken), ctx, token)
}

// GetAndDeleteByToken mocks base method.
func (m *MockRepo) GetAndDeleteByToken(ctx context.Context, token uuid.UUID) (entity.RefreshSession, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAndDeleteByToken", ctx, token)
	ret0, _ := ret[0].(entity.RefreshSession)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAndDeleteByToken indicates an expected call of GetAndDeleteByToken.
func (mr *MockRepoMockRecorder) GetAndDeleteByToken(ctx, token any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAndDeleteByToken", reflect.TypeOf((*MockRepo)(nil).GetAndDeleteByToken), ctx, token)
}
