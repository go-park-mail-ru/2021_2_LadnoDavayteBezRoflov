package repositories

import (
	"backendServer/app/api/models"
)

type CheckListItemRepository interface {
	Create(checkListItem *models.CheckListItem) (err error)
	Update(checkListItem *models.CheckListItem) (err error)
	Delete(chliid uint) (err error)
	GetByID(chliid uint) (checkListItem *models.CheckListItem, err error)
}
