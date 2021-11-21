package repositories

import (
	models2 "backendServer/app/api/models"
)

type CardRepository interface {
	Create(card *models2.Card) (err error)
	Update(card *models2.Card) (err error)
	Delete(cid uint) (err error)
	GetByID(cid uint) (card *models2.Card, err error)
	GetCardComments(cid uint) (comments *[]models2.Comment, err error)
	GetCardCheckLists(cid uint) (checkLists *[]models2.CheckList, err error)
	Move(fromPos, toPos, fromCardListID, toCardListID uint) (err error)
}
