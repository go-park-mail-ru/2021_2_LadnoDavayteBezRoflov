package usecases

import "backendServer/app/models"

type CheckListItemUseCase interface {
	CreateCheckListItem(checkListItem *models.CheckListItem) (chliid uint, err error)
	GetCheckListItem(uid, chliid uint) (checkListItem *models.CheckListItem, err error)
	UpdateCheckListItem(uid uint, checkListItem *models.CheckListItem) (err error)
	DeleteCheckListItem(uid, chliid uint) (err error)
}
