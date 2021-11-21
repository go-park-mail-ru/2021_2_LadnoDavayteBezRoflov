package repositories

import (
	models2 "backendServer/app/api/models"
)

type CardListRepository interface {
	Create(cardList *models2.CardList) (err error)
	Update(cardList *models2.CardList) (err error)
	Delete(clid uint) (err error)
	GetByID(clid uint) (cardList *models2.CardList, err error)
	GetCardListCards(clid uint) (cards *[]models2.Card, err error)
	Move(fromPos, toPos, bid uint) (err error)
}
