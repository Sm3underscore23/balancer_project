// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go
//
// Generated by this command:
//
//	mockgen -source=repository.go -destination=mocks/mock.go
//

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	model "balancer/internal/model"
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockLimitsRepository is a mock of LimitsRepository interface.
type MockLimitsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockLimitsRepositoryMockRecorder
	isgomock struct{}
}

// MockLimitsRepositoryMockRecorder is the mock recorder for MockLimitsRepository.
type MockLimitsRepositoryMockRecorder struct {
	mock *MockLimitsRepository
}

// NewMockLimitsRepository creates a new mock instance.
func NewMockLimitsRepository(ctrl *gomock.Controller) *MockLimitsRepository {
	mock := &MockLimitsRepository{ctrl: ctrl}
	mock.recorder = &MockLimitsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLimitsRepository) EXPECT() *MockLimitsRepositoryMockRecorder {
	return m.recorder
}

// CreateClientLimits mocks base method.
func (m *MockLimitsRepository) CreateClientLimits(ctx context.Context, info model.ClientLimits) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClientLimits", ctx, info)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateClientLimits indicates an expected call of CreateClientLimits.
func (mr *MockLimitsRepositoryMockRecorder) CreateClientLimits(ctx, info any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClientLimits", reflect.TypeOf((*MockLimitsRepository)(nil).CreateClientLimits), ctx, info)
}

// DeleteClientLimits mocks base method.
func (m *MockLimitsRepository) DeleteClientLimits(ctx context.Context, clientId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteClientLimits", ctx, clientId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteClientLimits indicates an expected call of DeleteClientLimits.
func (mr *MockLimitsRepositoryMockRecorder) DeleteClientLimits(ctx, clientId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteClientLimits", reflect.TypeOf((*MockLimitsRepository)(nil).DeleteClientLimits), ctx, clientId)
}

// GetClientLimits mocks base method.
func (m *MockLimitsRepository) GetClientLimits(ctx context.Context, clientId string) (model.ClientLimits, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClientLimits", ctx, clientId)
	ret0, _ := ret[0].(model.ClientLimits)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClientLimits indicates an expected call of GetClientLimits.
func (mr *MockLimitsRepositoryMockRecorder) GetClientLimits(ctx, clientId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClientLimits", reflect.TypeOf((*MockLimitsRepository)(nil).GetClientLimits), ctx, clientId)
}

// IsClientExists mocks base method.
func (m *MockLimitsRepository) IsClientExists(ctx context.Context, userId string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClientExists", ctx, userId)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsClientExists indicates an expected call of IsClientExists.
func (mr *MockLimitsRepositoryMockRecorder) IsClientExists(ctx, userId any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClientExists", reflect.TypeOf((*MockLimitsRepository)(nil).IsClientExists), ctx, userId)
}

// UpdateClientLimits mocks base method.
func (m *MockLimitsRepository) UpdateClientLimits(ctx context.Context, updateData model.ClientLimits) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateClientLimits", ctx, updateData)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateClientLimits indicates an expected call of UpdateClientLimits.
func (mr *MockLimitsRepositoryMockRecorder) UpdateClientLimits(ctx, updateData any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateClientLimits", reflect.TypeOf((*MockLimitsRepository)(nil).UpdateClientLimits), ctx, updateData)
}
