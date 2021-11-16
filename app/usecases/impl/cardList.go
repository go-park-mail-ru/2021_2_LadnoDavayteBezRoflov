package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	customErrors "backendServer/pkg/errors"
)

type CardListUseCaseImpl struct {
	cardListRepository repositories.CardListRepository
	userRepository     repositories.UserRepository
	teamRepository     repositories.TeamRepository
	boardRepository    repositories.BoardRepository
}

func CreateCardListUseCase(
	cardListRepository repositories.CardListRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
	boardRepository repositories.BoardRepository,
) usecases.CardListUseCase {
	return &CardListUseCaseImpl{
		cardListRepository: cardListRepository,
		userRepository:     userRepository,
		teamRepository:     teamRepository,
		boardRepository:    boardRepository,
	}
}

func (cardListUseCase *CardListUseCaseImpl) CreateCardList(cardList *models.CardList) (clid uint, err error) {
	err = cardListUseCase.cardListRepository.Create(cardList)
	if err != nil {
		return 0, err
	}
	return cardList.CLID, nil
}

func (cardListUseCase *CardListUseCaseImpl) GetCardList(uid, clid uint) (cardList *models.CardList, err error) {
	isAccessed, err := cardListUseCase.userRepository.IsCardListAccessed(uid, clid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

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
	isAccessed, err := cardListUseCase.userRepository.IsCardListAccessed(uid, cardList.CLID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return cardListUseCase.cardListRepository.Update(cardList)
}

func (cardListUseCase *CardListUseCaseImpl) DeleteCardList(uid, clid uint) (err error) {
	isAccessed, err := cardListUseCase.userRepository.IsCardListAccessed(uid, clid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return cardListUseCase.cardListRepository.Delete(clid)
}
