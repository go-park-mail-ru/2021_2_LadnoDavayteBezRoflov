package usecases

import (
	"backendServer/app/api/models"
)

type CardUseCase interface {
	CreateCard(card *models.Card) (cid uint, err error)
	GetCard(uid, cid uint) (card *models.Card, err error)
	UpdateCard(uid uint, card *models.Card) (err error)
	DeleteCard(uid, cid uint) (err error)
	ToggleUser(uid, cid, toggledUserID uint) (card *models.Card, err error)
	ToggleTag(uid, cid, toggledTagID uint) (card *models.Card, err error)
	UpdateAccessPath(uid, cid uint) (newAccessLink string, err error)
	AddUserViaLink(uid uint, accessPath string) (card *models.Card, err error)
}
