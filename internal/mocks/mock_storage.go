// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/tank4gun/go-musthave-diploma-tpl/internal/storage (interfaces: Storage)

// Package mocks is a generated GoMock package.
package mocks

import (
	"context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/tank4gun/go-musthave-diploma-tpl/internal/storage"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddOrderForUser mocks base method.
func (m *MockStorage) AddOrderForUser(ctx context.Context, arg0, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrderForUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOrderForUser indicates an expected call of AddOrderForUser.
func (mr *MockStorageMockRecorder) AddOrderForUser(ctx context.Context, arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrderForUser", reflect.TypeOf((*MockStorage)(nil).AddOrderForUser), arg0, arg1)
}

// AddWithdrawalForUser mocks base method.
func (m *MockStorage) AddWithdrawalForUser(ctx context.Context, arg0 string, arg1 storage.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWithdrawalForUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddWithdrawalForUser indicates an expected call of AddWithdrawalForUser.
func (mr *MockStorageMockRecorder) AddWithdrawalForUser(ctx context.Context, arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWithdrawalForUser", reflect.TypeOf((*MockStorage)(nil).AddWithdrawalForUser), arg0, arg1)
}

// GetOrdersByUser mocks base method.
func (m *MockStorage) GetOrdersByUser(ctx context.Context, arg0 string) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUser", arg0)
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUser indicates an expected call of GetOrdersByUser.
func (mr *MockStorageMockRecorder) GetOrdersByUser(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUser", reflect.TypeOf((*MockStorage)(nil).GetOrdersByUser), arg0)
}

// GetOrdersInProgress mocks base method.
func (m *MockStorage) GetOrdersInProgress(ctx context.Context) ([]storage.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersInProgress")
	ret0, _ := ret[0].([]storage.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersInProgress indicates an expected call of GetOrdersInProgress.
func (mr *MockStorageMockRecorder) GetOrdersInProgress(ctx context.Context) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersInProgress", reflect.TypeOf((*MockStorage)(nil).GetOrdersInProgress))
}

// GetUserBalance mocks base method.
func (m *MockStorage) GetUserBalance(ctx context.Context, arg0 string) (storage.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalance", arg0)
	ret0, _ := ret[0].(storage.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalance indicates an expected call of GetUserBalance.
func (mr *MockStorageMockRecorder) GetUserBalance(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalance", reflect.TypeOf((*MockStorage)(nil).GetUserBalance), arg0)
}

// GetUserByLogin mocks base method.
func (m *MockStorage) GetUserByLogin(ctx context.Context, arg0 storage.Auth) (storage.Auth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0)
	ret0, _ := ret[0].(storage.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStorageMockRecorder) GetUserByLogin(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorage)(nil).GetUserByLogin), arg0)
}

// GetWithdrawalsForUser mocks base method.
func (m *MockStorage) GetWithdrawalsForUser(ctx context.Context, arg0 string) ([]storage.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsForUser", arg0)
	ret0, _ := ret[0].([]storage.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsForUser indicates an expected call of GetWithdrawalsForUser.
func (mr *MockStorageMockRecorder) GetWithdrawalsForUser(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsForUser", reflect.TypeOf((*MockStorage)(nil).GetWithdrawalsForUser), arg0)
}

// Register mocks base method.
func (m *MockStorage) Register(ctx context.Context, arg0 storage.Auth) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockStorageMockRecorder) Register(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockStorage)(nil).Register), arg0)
}

// UpdateOrder mocks base method.
func (m *MockStorage) UpdateOrder(ctx context.Context, arg0 storage.OrderFromBlackBox) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStorageMockRecorder) UpdateOrder(ctx context.Context, arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStorage)(nil).UpdateOrder), arg0)
}
