// Code generated by MockGen. DO NOT EDIT.
// Source: param.go

// Package mock_node is a generated GoMock package.
package mock_node

import (
	reflect "reflect"

	node "git.in.zhihu.com/antispam/datasupply/node"
	gomock "github.com/golang/mock/gomock"
)

// MockINodeParams is a mock of INodeParams interface.
type MockINodeParams struct {
	ctrl     *gomock.Controller
	recorder *MockINodeParamsMockRecorder
}

// MockINodeParamsMockRecorder is the mock recorder for MockINodeParams.
type MockINodeParamsMockRecorder struct {
	mock *MockINodeParams
}

// NewMockINodeParams creates a new mock instance.
func NewMockINodeParams(ctrl *gomock.Controller) *MockINodeParams {
	mock := &MockINodeParams{ctrl: ctrl}
	mock.recorder = &MockINodeParamsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockINodeParams) EXPECT() *MockINodeParamsMockRecorder {
	return m.recorder
}

// AddFuncParams mocks base method.
func (m *MockINodeParams) AddFuncParams(funcName string, funcParams []node.Param) node.INodeParams {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddFuncParams", funcName, funcParams)
	ret0, _ := ret[0].(node.INodeParams)
	return ret0
}

// AddFuncParams indicates an expected call of AddFuncParams.
func (mr *MockINodeParamsMockRecorder) AddFuncParams(funcName, funcParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddFuncParams", reflect.TypeOf((*MockINodeParams)(nil).AddFuncParams), funcName, funcParams)
}

// GetParamMap mocks base method.
func (m *MockINodeParams) GetParamMap() map[string][]node.Param {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParamMap")
	ret0, _ := ret[0].(map[string][]node.Param)
	return ret0
}

// GetParamMap indicates an expected call of GetParamMap.
func (mr *MockINodeParamsMockRecorder) GetParamMap() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParamMap", reflect.TypeOf((*MockINodeParams)(nil).GetParamMap))
}

// GetParamVariables mocks base method.
func (m *MockINodeParams) GetParamVariables() []node.Param {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParamVariables")
	ret0, _ := ret[0].([]node.Param)
	return ret0
}

// GetParamVariables indicates an expected call of GetParamVariables.
func (mr *MockINodeParamsMockRecorder) GetParamVariables() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParamVariables", reflect.TypeOf((*MockINodeParams)(nil).GetParamVariables))
}

// GetParams mocks base method.
func (m *MockINodeParams) GetParams() []node.Param {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParams")
	ret0, _ := ret[0].([]node.Param)
	return ret0
}

// GetParams indicates an expected call of GetParams.
func (mr *MockINodeParamsMockRecorder) GetParams() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParams", reflect.TypeOf((*MockINodeParams)(nil).GetParams))
}

// GetParamsByFunc mocks base method.
func (m *MockINodeParams) GetParamsByFunc(funcName string) []node.Param {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetParamsByFunc", funcName)
	ret0, _ := ret[0].([]node.Param)
	return ret0
}

// GetParamsByFunc indicates an expected call of GetParamsByFunc.
func (mr *MockINodeParamsMockRecorder) GetParamsByFunc(funcName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetParamsByFunc", reflect.TypeOf((*MockINodeParams)(nil).GetParamsByFunc), funcName)
}
