// Code generated by MockGen. DO NOT EDIT.
// Source: cardList.go

// Package mocks is a generated GoMock package.
package mocks

import (
	models "backendServer/app/api/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockCardListRepository is a mock of CardListRepository interface.
type MockCardListRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCardListRepositoryMockRecorder
}

// MockCardListRepositoryMockRecorder is the mock recorder for MockCardListRepository.
type MockCardListRepositoryMockRecorder struct {
	mock *MockCardListRepository
}

// NewMockCardListRepository creates a new mock instance.
func NewMockCardListRepository(ctrl *gomock.Controller) *MockCardListRepository {
	mock := &MockCardListRepository{ctrl: ctrl}
	mock.recorder = &MockCardListRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCardListRepository) EXPECT() *MockCardListRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCardListRepository) Create(cardList *models.CardList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", cardList)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockCardListRepositoryMockRecorder) Create(cardList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCardListRepository)(nil).Create), cardList)
}

// Delete mocks base method.
func (m *MockCardListRepository) Delete(clid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", clid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCardListRepositoryMockRecorder) Delete(clid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCardListRepository)(nil).Delete), clid)
}

// GetByID mocks base method.
func (m *MockCardListRepository) GetByID(clid uint) (*models.CardList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", clid)
	ret0, _ := ret[0].(*models.CardList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCardListRepositoryMockRecorder) GetByID(clid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCardListRepository)(nil).GetByID), clid)
}

// GetCardListCards mocks base method.
func (m *MockCardListRepository) GetCardListCards(clid uint) (*[]models.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCardListCards", clid)
	ret0, _ := ret[0].(*[]models.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCardListCards indicates an expected call of GetCardListCards.
func (mr *MockCardListRepositoryMockRecorder) GetCardListCards(clid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCardListCards", reflect.TypeOf((*MockCardListRepository)(nil).GetCardListCards), clid)
}

// Move mocks base method.
func (m *MockCardListRepository) Move(fromPos, toPos, bid uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Move", fromPos, toPos, bid)
	ret0, _ := ret[0].(error)
	return ret0
}

// Move indicates an expected call of Move.
func (mr *MockCardListRepositoryMockRecorder) Move(fromPos, toPos, bid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Move", reflect.TypeOf((*MockCardListRepository)(nil).Move), fromPos, toPos, bid)
}

// Update mocks base method.
func (m *MockCardListRepository) Update(cardList *models.CardList) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", cardList)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockCardListRepositoryMockRecorder) Update(cardList interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCardListRepository)(nil).Update), cardList)
}