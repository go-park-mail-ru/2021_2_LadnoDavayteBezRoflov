// Code generated by MockGen. DO NOT EDIT.
// Source: checkListItem.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCheckListItemRepository is a mock of CheckListItemRepository interface.
type MockCheckListItemRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCheckListItemRepositoryMockRecorder
}

// MockCheckListItemRepositoryMockRecorder is the mock recorder for MockCheckListItemRepository.
type MockCheckListItemRepositoryMockRecorder struct {
	mock *MockCheckListItemRepository
}

// NewMockCheckListItemRepository creates a new mock instance.
func NewMockCheckListItemRepository(ctrl *gomock.Controller) *MockCheckListItemRepository {
	mock := &MockCheckListItemRepository{ctrl: ctrl}
	mock.recorder = &MockCheckListItemRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCheckListItemRepository) EXPECT() *MockCheckListItemRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCheckListItemRepository) Create(checkListItem *models.CheckListItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", checkListItem)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCheckListItemRepositoryMockRecorder) Create(checkListItem interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCheckListItemRepository)(nil).Create), checkListItem)
}

// Delete mocks base method.
func (m *MockCheckListItemRepository) Delete(chliid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", chliid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCheckListItemRepositoryMockRecorder) Delete(chliid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCheckListItemRepository)(nil).Delete), chliid)
}

// GetByID mocks base method.
func (m *MockCheckListItemRepository) GetByID(chliid uint) (*models.CheckListItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", chliid)
	ret0, _ := ret[0].(*models.CheckListItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCheckListItemRepositoryMockRecorder) GetByID(chliid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCheckListItemRepository)(nil).GetByID), chliid)
}

// Update mocks base method.
func (m *MockCheckListItemRepository) Update(checkListItem *models.CheckListItem) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", checkListItem)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCheckListItemRepositoryMockRecorder) Update(checkListItem interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCheckListItemRepository)(nil).Update), checkListItem)
}
