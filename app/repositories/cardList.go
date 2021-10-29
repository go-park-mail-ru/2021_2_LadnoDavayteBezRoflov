package repositories

import "backendServer/app/models"

type CardListRepository interface { // TODO реализовать перенос в рамках доски
	Create(cardList *models.CardList) (err error)
	Update(cardList *models.CardList) (err error)
	Delete(clid uint) (err error)
	GetByID(clid uint) (cardList *models.CardList, err error)
	GetCardListCards(clid uint) (cards *[]models.Card, err error)
}
