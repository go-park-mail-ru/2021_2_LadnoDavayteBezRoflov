// Code generated by MockGen. DO NOT EDIT.
// Source: userSearch.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserSearchUseCase is a mock of UserSearchUseCase interface.
type MockUserSearchUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUserSearchUseCaseMockRecorder
}

// MockUserSearchUseCaseMockRecorder is the mock recorder for MockUserSearchUseCase.
type MockUserSearchUseCaseMockRecorder struct {
	mock *MockUserSearchUseCase
}

// NewMockUserSearchUseCase creates a new mock instance.
func NewMockUserSearchUseCase(ctrl *gomock.Controller) *MockUserSearchUseCase {
	mock := &MockUserSearchUseCase{ctrl: ctrl}
	mock.recorder = &MockUserSearchUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserSearchUseCase) EXPECT() *MockUserSearchUseCaseMockRecorder {
	return m.recorder
}

// FindForBoard mocks base method.
func (m *MockUserSearchUseCase) FindForBoard(uid, bid uint, text string) (*[]models.UserSearchInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindForBoard", uid, bid, text)
	ret0, _ := ret[0].(*[]models.UserSearchInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindForBoard indicates an expected call of FindForBoard.
func (mr *MockUserSearchUseCaseMockRecorder) FindForBoard(uid, bid, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindForBoard", reflect.TypeOf((*MockUserSearchUseCase)(nil).FindForBoard), uid, bid, text)
}

// FindForCard mocks base method.
func (m *MockUserSearchUseCase) FindForCard(uid, cid uint, text string) (*[]models.UserSearchInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindForCard", uid, cid, text)
	ret0, _ := ret[0].(*[]models.UserSearchInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindForCard indicates an expected call of FindForCard.
func (mr *MockUserSearchUseCaseMockRecorder) FindForCard(uid, cid, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindForCard", reflect.TypeOf((*MockUserSearchUseCase)(nil).FindForCard), uid, cid, text)
}

// FindForTeam mocks base method.
func (m *MockUserSearchUseCase) FindForTeam(uid, tid uint, text string) (*[]models.UserSearchInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindForTeam", uid, tid, text)
	ret0, _ := ret[0].(*[]models.UserSearchInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindForTeam indicates an expected call of FindForTeam.
func (mr *MockUserSearchUseCaseMockRecorder) FindForTeam(uid, tid, text interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindForTeam", reflect.TypeOf((*MockUserSearchUseCase)(nil).FindForTeam), uid, tid, text)
}
