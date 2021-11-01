package repositories

import "backendServer/app/models"

type CardRepository interface {
	Create(card *models.Card) (err error)
	Update(card *models.Card) (err error)
	Delete(cid uint) (err error)
	GetByID(cid uint) (card *models.Card, err error)
	Move(fromPos, toPos, fromCardListID, toCardListID uint) (err error)
}
