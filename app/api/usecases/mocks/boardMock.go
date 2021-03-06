// Code generated by MockGen. DO NOT EDIT.
// Source: board.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBoardUseCase is a mock of BoardUseCase interface.
type MockBoardUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockBoardUseCaseMockRecorder
}

// MockBoardUseCaseMockRecorder is the mock recorder for MockBoardUseCase.
type MockBoardUseCaseMockRecorder struct {
	mock *MockBoardUseCase
}

// NewMockBoardUseCase creates a new mock instance.
func NewMockBoardUseCase(ctrl *gomock.Controller) *MockBoardUseCase {
	mock := &MockBoardUseCase{ctrl: ctrl}
	mock.recorder = &MockBoardUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBoardUseCase) EXPECT() *MockBoardUseCaseMockRecorder {
	return m.recorder
}

// AddUserViaLink mocks base method.
func (m *MockBoardUseCase) AddUserViaLink(uid uint, accessPath string) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserViaLink", uid, accessPath)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUserViaLink indicates an expected call of AddUserViaLink.
func (mr *MockBoardUseCaseMockRecorder) AddUserViaLink(uid, accessPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserViaLink", reflect.TypeOf((*MockBoardUseCase)(nil).AddUserViaLink), uid, accessPath)
}

// CreateBoard mocks base method.
func (m *MockBoardUseCase) CreateBoard(board *models.Board) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateBoard", board)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateBoard indicates an expected call of CreateBoard.
func (mr *MockBoardUseCaseMockRecorder) CreateBoard(board interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateBoard", reflect.TypeOf((*MockBoardUseCase)(nil).CreateBoard), board)
}

// DeleteBoard mocks base method.
func (m *MockBoardUseCase) DeleteBoard(uid, bid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBoard", uid, bid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBoard indicates an expected call of DeleteBoard.
func (mr *MockBoardUseCaseMockRecorder) DeleteBoard(uid, bid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBoard", reflect.TypeOf((*MockBoardUseCase)(nil).DeleteBoard), uid, bid)
}

// GetBoard mocks base method.
func (m *MockBoardUseCase) GetBoard(uid, bid uint) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBoard", uid, bid)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBoard indicates an expected call of GetBoard.
func (mr *MockBoardUseCaseMockRecorder) GetBoard(uid, bid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBoard", reflect.TypeOf((*MockBoardUseCase)(nil).GetBoard), uid, bid)
}

// GetUserBoards mocks base method.
func (m *MockBoardUseCase) GetUserBoards(uid uint) (*[]models.Team, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBoards", uid)
	ret0, _ := ret[0].(*[]models.Team)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBoards indicates an expected call of GetUserBoards.
func (mr *MockBoardUseCaseMockRecorder) GetUserBoards(uid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBoards", reflect.TypeOf((*MockBoardUseCase)(nil).GetUserBoards), uid)
}

// ToggleUser mocks base method.
func (m *MockBoardUseCase) ToggleUser(uid, bid, toggledUserID uint) (*models.Board, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToggleUser", uid, bid, toggledUserID)
	ret0, _ := ret[0].(*models.Board)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToggleUser indicates an expected call of ToggleUser.
func (mr *MockBoardUseCaseMockRecorder) ToggleUser(uid, bid, toggledUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleUser", reflect.TypeOf((*MockBoardUseCase)(nil).ToggleUser), uid, bid, toggledUserID)
}

// UpdateAccessPath mocks base method.
func (m *MockBoardUseCase) UpdateAccessPath(uid, bid uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccessPath", uid, bid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccessPath indicates an expected call of UpdateAccessPath.
func (mr *MockBoardUseCaseMockRecorder) UpdateAccessPath(uid, bid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccessPath", reflect.TypeOf((*MockBoardUseCase)(nil).UpdateAccessPath), uid, bid)
}

// UpdateBoard mocks base method.
func (m *MockBoardUseCase) UpdateBoard(uid uint, board *models.Board) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBoard", uid, board)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateBoard indicates an expected call of UpdateBoard.
func (mr *MockBoardUseCaseMockRecorder) UpdateBoard(uid, board interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBoard", reflect.TypeOf((*MockBoardUseCase)(nil).UpdateBoard), uid, board)
}
