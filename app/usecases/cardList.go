package usecases

import "backendServer/app/models"

type CardListUseCase interface {
	CreateCardList(cardList *models.CardList) (clid uint, err error)
	GetCardList(uid, clid uint) (cardList *models.CardList, err error)
	UpdateCardList(uid uint, cardList *models.CardList) (err error)
	DeleteCardList(uid, clid uint) (err error)
}
