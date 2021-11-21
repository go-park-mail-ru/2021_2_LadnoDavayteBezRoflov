package repositories

import (
	models2 "backendServer/app/api/models"
)

type CheckListRepository interface {
	Create(checkList *models2.CheckList) (err error)
	Update(checkList *models2.CheckList) (err error)
	Delete(chlid uint) (err error)
	GetByID(chlid uint) (checkList *models2.CheckList, err error)
	GetCheckListItems(chlid uint) (checkListItems *[]models2.CheckListItem, err error)
}
