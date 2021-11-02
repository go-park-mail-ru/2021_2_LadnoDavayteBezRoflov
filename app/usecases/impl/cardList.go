package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
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

//func (cardListUseCase *CardListUseCaseImpl) isUserHaveAccessToCardList(uid, clid uint) (isAccessed bool, err error) {
//	teams, err := cardListUseCase.userRepository.GetUserTeams(uid)
//	if err != nil {
//		return
//	}
//
//	for _, team := range *teams {
//		boards, boardsErr := boardUseCase.teamRepository.GetTeamBoards(team.TID)
//		if boardsErr != nil {
//			err = boardsErr
//			return
//		}
//		for _, board := range *boards {
//			if board.BID == bid {
//				isAccessed = true
//				return
//			}
//		}
//	}
//	isAccessed = false
//	return
//}
