package usecases

import (
	"backendServer/app/api/models"
)

type CheckListUseCase interface {
	CreateCheckList(checkList *models.CheckList) (chlid uint, err error)
	GetCheckList(uid, chlid uint) (checkList *models.CheckList, err error)
	UpdateCheckList(uid uint, checkList *models.CheckList) (err error)
	DeleteCheckList(uid, chlid uint) (err error)
}
