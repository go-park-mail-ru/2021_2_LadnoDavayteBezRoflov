// Code generated by MockGen. DO NOT EDIT.
// Source: session.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockSessionUseCase is a mock of SessionUseCase interface.
type MockSessionUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockSessionUseCaseMockRecorder
}

// MockSessionUseCaseMockRecorder is the mock recorder for MockSessionUseCase.
type MockSessionUseCaseMockRecorder struct {
	mock *MockSessionUseCase
}

// NewMockSessionUseCase creates a new mock instance.
func NewMockSessionUseCase(ctrl *gomock.Controller) *MockSessionUseCase {
	mock := &MockSessionUseCase{ctrl: ctrl}
	mock.recorder = &MockSessionUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSessionUseCase) EXPECT() *MockSessionUseCaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockSessionUseCase) Create(user *models.User) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", user)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockSessionUseCaseMockRecorder) Create(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockSessionUseCase)(nil).Create), user)
}

// Delete mocks base method.
func (m *MockSessionUseCase) Delete(sid string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", sid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockSessionUseCaseMockRecorder) Delete(sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockSessionUseCase)(nil).Delete), sid)
}

// Get mocks base method.
func (m *MockSessionUseCase) Get(sid string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", sid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockSessionUseCaseMockRecorder) Get(sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockSessionUseCase)(nil).Get), sid)
}

// GetUID mocks base method.
func (m *MockSessionUseCase) GetUID(sid string) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUID", sid)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUID indicates an expected call of GetUID.
func (mr *MockSessionUseCaseMockRecorder) GetUID(sid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUID", reflect.TypeOf((*MockSessionUseCase)(nil).GetUID), sid)
}