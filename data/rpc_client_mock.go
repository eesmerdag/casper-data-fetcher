// Code generated by MockGen. DO NOT EDIT.
// Source: data/rpc_client.go

// Package mock_data is a generated GoMock package.
package data

import (
	reflect "reflect"

	sdk "github.com/casper-ecosystem/casper-golang-sdk/sdk"
	gomock "github.com/golang/mock/gomock"
)

// MockRpcClientI is a mock of RpcClientI interface.
type MockRpcClientI struct {
	ctrl     *gomock.Controller
	recorder *MockRpcClientIMockRecorder
}

// MockRpcClientIMockRecorder is the mock recorder for MockRpcClientI.
type MockRpcClientIMockRecorder struct {
	mock *MockRpcClientI
}

// NewMockRpcClientI creates a new mock instance.
func NewMockRpcClientI(ctrl *gomock.Controller) *MockRpcClientI {
	mock := &MockRpcClientI{ctrl: ctrl}
	mock.recorder = &MockRpcClientIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRpcClientI) EXPECT() *MockRpcClientIMockRecorder {
	return m.recorder
}

// GetBlockByHeight mocks base method.
func (m *MockRpcClientI) GetBlockByHeight(height uint64) (sdk.BlockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockByHeight", height)
	ret0, _ := ret[0].(sdk.BlockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockByHeight indicates an expected call of GetBlockByHeight.
func (mr *MockRpcClientIMockRecorder) GetBlockByHeight(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockByHeight", reflect.TypeOf((*MockRpcClientI)(nil).GetBlockByHeight), height)
}

// GetBlockTransfersByHeight mocks base method.
func (m *MockRpcClientI) GetBlockTransfersByHeight(height uint64) ([]sdk.TransferResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockTransfersByHeight", height)
	ret0, _ := ret[0].([]sdk.TransferResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockTransfersByHeight indicates an expected call of GetBlockTransfersByHeight.
func (mr *MockRpcClientIMockRecorder) GetBlockTransfersByHeight(height interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockTransfersByHeight", reflect.TypeOf((*MockRpcClientI)(nil).GetBlockTransfersByHeight), height)
}

// GetLatestBlock mocks base method.
func (m *MockRpcClientI) GetLatestBlock() (sdk.BlockResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestBlock")
	ret0, _ := ret[0].(sdk.BlockResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLatestBlock indicates an expected call of GetLatestBlock.
func (mr *MockRpcClientIMockRecorder) GetLatestBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestBlock", reflect.TypeOf((*MockRpcClientI)(nil).GetLatestBlock))
}
