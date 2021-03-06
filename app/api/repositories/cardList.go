package repositories

import (
	"backendServer/app/api/models"
)

type CardListRepository interface {
	Create(cardList *models.CardList) (err error)
	Update(cardList *models.CardList) (err error)
	Delete(clid uint) (err error)
	GetByID(clid uint) (cardList *models.CardList, err error)
	GetCardListCards(clid uint) (cards *[]models.Card, err error)
	Move(fromPos, toPos, bid uint) (err error)
}
