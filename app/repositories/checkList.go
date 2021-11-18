package repositories

import "backendServer/app/models"

type CheckListRepository interface {
	Create(checkList *models.CheckList) (err error)
	Update(checkList *models.CheckList) (err error)
	Delete(chlid uint) (err error)
	GetByID(chlid uint) (checkList *models.CheckList, err error)
	GetCheckListItems(chlid uint) (checkListItems *[]models.CheckListItem, err error)
}
