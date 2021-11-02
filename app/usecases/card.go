package usecases

import "backendServer/app/models"

type CardUseCase interface {
	CreateCard(card *models.Card) (cid uint, err error)
	GetCard(uid, cid uint) (card *models.Card, err error)
	UpdateCard(uid uint, card *models.Card) (err error)
	DeleteCard(uid, cid uint) (err error)
}
