package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories/mocks"
	customErrors "backendServer/pkg/errors"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func createCardListRepoMocks(controller *gomock.Controller) (*mocks.MockCardListRepository, *mocks.MockUserRepository) {
	cardListRepoMock := mocks.NewMockCardListRepository(controller)
	userRepoMock := mocks.NewMockUserRepository(controller)
	return cardListRepoMock, userRepoMock
}

func TestCreateCardList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardListRepoMock, userRepoMock := createCardListRepoMocks(ctrl)
	cardListUseCase := CreateCardListUseCase(cardListRepoMock, userRepoMock)

	testCardList := new(models.CardList)
	err := faker.FakeData(testCardList)
	assert.NoError(t, err)

	// good
	cardListRepoMock.EXPECT().Create(testCardList).Return(nil)
	resClid, err := cardListUseCase.CreateCardList(testCardList)
	assert.NoError(t, err)
	assert.Equal(t, testCardList.CLID, resClid)

	// error can't create
	cardListRepoMock.EXPECT().Create(testCardList).Return(customErrors.ErrInternal)
	resClid, err = cardListUseCase.CreateCardList(testCardList)
	assert.Equal(t, uint(0), resClid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestGetCardList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardListRepoMock, userRepoMock := createCardListRepoMocks(ctrl)
	cardListUseCase := CreateCardListUseCase(cardListRepoMock, userRepoMock)

	uid := uint(1)
	clid := uint(1)

	testCardList := new(models.CardList)
	err := faker.FakeData(testCardList)
	assert.NoError(t, err)
	testCards := make([]models.Card, 3)
	for i := range testCards {
		err = faker.FakeData(&testCards[i])
		assert.NoError(t, err)
	}
	testCardList.Cards = testCards

	// success
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(true, nil)
	cardListRepoMock.EXPECT().GetByID(clid).Return(testCardList, nil)
	cardListRepoMock.EXPECT().GetCardListCards(clid).Return(&testCards, nil)
	resCardList, err := cardListUseCase.GetCardList(uid, clid)
	assert.NoError(t, err)
	assert.Equal(t, testCardList, resCardList)

	// error while checking access
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(false, customErrors.ErrInternal)
	_, err = cardListUseCase.GetCardList(uid, clid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(false, nil)
	_, err = cardListUseCase.GetCardList(uid, clid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't found
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(true, nil)
	cardListRepoMock.EXPECT().GetByID(clid).Return(nil, customErrors.ErrCardListNotFound)
	_, err = cardListUseCase.GetCardList(uid, clid)
	assert.Equal(t, customErrors.ErrCardListNotFound, err)

	// can't get check list items
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(true, nil)
	cardListRepoMock.EXPECT().GetByID(clid).Return(testCardList, nil)
	cardListRepoMock.EXPECT().GetCardListCards(clid).Return(nil, customErrors.ErrInternal)
	_, err = cardListUseCase.GetCardList(uid, clid)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestUpdateCardList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardListRepoMock, userRepoMock := createCardListRepoMocks(ctrl)
	cardListUseCase := CreateCardListUseCase(cardListRepoMock, userRepoMock)

	uid := uint(1)
	testCardList := new(models.CardList)
	err := faker.FakeData(testCardList)
	assert.NoError(t, err)

	// success
	userRepoMock.EXPECT().IsCardListAccessed(uid, testCardList.CLID).Return(true, nil)
	cardListRepoMock.EXPECT().Update(testCardList).Return(nil)
	err = cardListUseCase.UpdateCardList(uid, testCardList)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCardListAccessed(uid, testCardList.CLID).Return(false, customErrors.ErrInternal)
	err = cardListUseCase.UpdateCardList(uid, testCardList)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardListAccessed(uid, testCardList.CLID).Return(false, nil)
	err = cardListUseCase.UpdateCardList(uid, testCardList)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't update
	userRepoMock.EXPECT().IsCardListAccessed(uid, testCardList.CLID).Return(true, nil)
	cardListRepoMock.EXPECT().Update(testCardList).Return(customErrors.ErrInternal)
	err = cardListUseCase.UpdateCardList(uid, testCardList)
	assert.Equal(t, customErrors.ErrInternal, err)
}

func TestDeleteCardList(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	cardListRepoMock, userRepoMock := createCardListRepoMocks(ctrl)
	cardListUseCase := CreateCardListUseCase(cardListRepoMock, userRepoMock)

	uid := uint(1)
	clid := uint(1)

	// good
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(true, nil)
	cardListRepoMock.EXPECT().Delete(clid).Return(nil)
	err := cardListUseCase.DeleteCardList(uid, clid)
	assert.NoError(t, err)

	// error while checking access
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(false, customErrors.ErrInternal)
	err = cardListUseCase.DeleteCardList(uid, clid)
	assert.Equal(t, customErrors.ErrInternal, err)

	// no access
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(false, nil)
	err = cardListUseCase.DeleteCardList(uid, clid)
	assert.Equal(t, customErrors.ErrNoAccess, err)

	// can't delete
	userRepoMock.EXPECT().IsCardListAccessed(uid, clid).Return(true, nil)
	cardListRepoMock.EXPECT().Delete(clid).Return(customErrors.ErrInternal)
	err = cardListUseCase.DeleteCardList(uid, clid)
	assert.Equal(t, customErrors.ErrInternal, err)
}
