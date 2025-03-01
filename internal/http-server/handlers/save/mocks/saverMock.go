// Code generated by MockGen. DO NOT EDIT.
// Source: save.go

// Package mock_save is a generated GoMock package.
package mock_save

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockWalletSaver is a mock of WalletSaver interface.
type MockWalletSaver struct {
	ctrl     *gomock.Controller
	recorder *MockWalletSaverMockRecorder
}

// MockWalletSaverMockRecorder is the mock recorder for MockWalletSaver.
type MockWalletSaverMockRecorder struct {
	mock *MockWalletSaver
}

// NewMockWalletSaver creates a new mock instance.
func NewMockWalletSaver(ctrl *gomock.Controller) *MockWalletSaver {
	mock := &MockWalletSaver{ctrl: ctrl}
	mock.recorder = &MockWalletSaverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWalletSaver) EXPECT() *MockWalletSaverMockRecorder {
	return m.recorder
}

// SaveWallet mocks base method.
func (m *MockWalletSaver) SaveWallet(ctx context.Context, amount int) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWallet", ctx, amount)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveWallet indicates an expected call of SaveWallet.
func (mr *MockWalletSaverMockRecorder) SaveWallet(ctx, amount interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWallet", reflect.TypeOf((*MockWalletSaver)(nil).SaveWallet), ctx, amount)
}
