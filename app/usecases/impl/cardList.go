package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
)

type CardListUseCaseImpl struct {
	cardListRepository repositories.CardListRepository
}

func CreateCardListUseCase(cardListRepository repositories.CardListRepository) usecases.CardListUseCase {
	return &CardListUseCaseImpl{cardListRepository: cardListRepository}
}

func (cardListUseCase *CardListUseCaseImpl) CreateCardList(cardList *models.CardList) (clid uint, err error) {
	err = cardListUseCase.cardListRepository.Create(cardList)
	if err != nil {
		return 0, err
	}
	return cardList.CLID, nil
}

func (cardListUseCase *CardListUseCaseImpl) GetCardList(uid, clid uint) (cardList *models.CardList, err error) {
	// TODO добавить проверку на то, что пользователь имеет доступ к этому списку карт
	cardList, err = cardListUseCase.cardListRepository.GetByID(clid)
	if err != nil {
		return
	}

	cards, err := cardListUseCase.cardListRepository.GetCardListCards(clid)
	if err != nil {
		return
	}

	cardList.Cards = *cards
	return
}

func (cardListUseCase *CardListUseCaseImpl) UpdateCardList(uid uint, cardList *models.CardList) (err error) {
	// TODO добавить проверку на то, что пользователь имеет доступ к этому списку карт

	return cardListUseCase.cardListRepository.Update(cardList)
}

func (cardListUseCase *CardListUseCaseImpl) DeleteCardList(uid, clid uint) (err error) {
	// TODO добавить проверку на то, что пользователь имеет доступ к этому списку карт

	return cardListUseCase.cardListRepository.Delete(clid)
}
