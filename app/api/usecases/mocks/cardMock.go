// Code generated by MockGen. DO NOT EDIT.
// Source: card.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCardUseCase is a mock of CardUseCase interface.
type MockCardUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockCardUseCaseMockRecorder
}

// MockCardUseCaseMockRecorder is the mock recorder for MockCardUseCase.
type MockCardUseCaseMockRecorder struct {
	mock *MockCardUseCase
}

// NewMockCardUseCase creates a new mock instance.
func NewMockCardUseCase(ctrl *gomock.Controller) *MockCardUseCase {
	mock := &MockCardUseCase{ctrl: ctrl}
	mock.recorder = &MockCardUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardUseCase) EXPECT() *MockCardUseCaseMockRecorder {
	return m.recorder
}

// AddUserViaLink mocks base method.
func (m *MockCardUseCase) AddUserViaLink(uid uint, accessPath string) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUserViaLink", uid, accessPath)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUserViaLink indicates an expected call of AddUserViaLink.
func (mr *MockCardUseCaseMockRecorder) AddUserViaLink(uid, accessPath interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUserViaLink", reflect.TypeOf((*MockCardUseCase)(nil).AddUserViaLink), uid, accessPath)
}

// CreateCard mocks base method.
func (m *MockCardUseCase) CreateCard(card *models.Card) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateCard", card)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateCard indicates an expected call of CreateCard.
func (mr *MockCardUseCaseMockRecorder) CreateCard(card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateCard", reflect.TypeOf((*MockCardUseCase)(nil).CreateCard), card)
}

// DeleteCard mocks base method.
func (m *MockCardUseCase) DeleteCard(uid, cid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteCard", uid, cid)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteCard indicates an expected call of DeleteCard.
func (mr *MockCardUseCaseMockRecorder) DeleteCard(uid, cid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteCard", reflect.TypeOf((*MockCardUseCase)(nil).DeleteCard), uid, cid)
}

// GetCard mocks base method.
func (m *MockCardUseCase) GetCard(uid, cid uint) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCard", uid, cid)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCard indicates an expected call of GetCard.
func (mr *MockCardUseCaseMockRecorder) GetCard(uid, cid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCard", reflect.TypeOf((*MockCardUseCase)(nil).GetCard), uid, cid)
}

// ToggleTag mocks base method.
func (m *MockCardUseCase) ToggleTag(uid, cid, toggledTagID uint) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToggleTag", uid, cid, toggledTagID)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToggleTag indicates an expected call of ToggleTag.
func (mr *MockCardUseCaseMockRecorder) ToggleTag(uid, cid, toggledTagID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleTag", reflect.TypeOf((*MockCardUseCase)(nil).ToggleTag), uid, cid, toggledTagID)
}

// ToggleUser mocks base method.
func (m *MockCardUseCase) ToggleUser(uid, cid, toggledUserID uint) (*models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToggleUser", uid, cid, toggledUserID)
	ret0, _ := ret[0].(*models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ToggleUser indicates an expected call of ToggleUser.
func (mr *MockCardUseCaseMockRecorder) ToggleUser(uid, cid, toggledUserID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToggleUser", reflect.TypeOf((*MockCardUseCase)(nil).ToggleUser), uid, cid, toggledUserID)
}

// UpdateAccessPath mocks base method.
func (m *MockCardUseCase) UpdateAccessPath(uid, cid uint) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccessPath", uid, cid)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateAccessPath indicates an expected call of UpdateAccessPath.
func (mr *MockCardUseCaseMockRecorder) UpdateAccessPath(uid, cid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccessPath", reflect.TypeOf((*MockCardUseCase)(nil).UpdateAccessPath), uid, cid)
}

// UpdateCard mocks base method.
func (m *MockCardUseCase) UpdateCard(uid uint, card *models.Card) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCard", uid, card)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCard indicates an expected call of UpdateCard.
func (mr *MockCardUseCaseMockRecorder) UpdateCard(uid, card interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCard", reflect.TypeOf((*MockCardUseCase)(nil).UpdateCard), uid, card)
}
