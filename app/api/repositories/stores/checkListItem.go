package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	"gorm.io/gorm"
)

type CheckListItemStore struct {
	db *gorm.DB
}

func CreateCheckListItemRepository(db *gorm.DB) repositories.CheckListItemRepository {
	return &CheckListItemStore{db: db}
}

func (checkListItemStore *CheckListItemStore) Create(checkListItem *models.CheckListItem) (err error) {
	return checkListItemStore.db.Create(checkListItem).Error
}

func (checkListItemStore *CheckListItemStore) Update(checkListItem *models.CheckListItem) (err error) {
	oldCheckListItem, err := checkListItemStore.GetByID(checkListItem.CHLIID)
	if err != nil {
		return
	}

	if checkListItem.Text != "" && checkListItem.Text != oldCheckListItem.Text {
		oldCheckListItem.Text = checkListItem.Text
	}

	if checkListItem.Status != oldCheckListItem.Status {
		oldCheckListItem.Status = checkListItem.Status
	}

	return checkListItemStore.db.Save(oldCheckListItem).Error
}

func (checkListItemStore *CheckListItemStore) Delete(chliid uint) (err error) {
	return checkListItemStore.db.Delete(&models.CheckListItem{}, chliid).Error
}

func (checkListItemStore *CheckListItemStore) GetByID(chliid uint) (*models.CheckListItem, error) {
	checkListItem := new(models.CheckListItem)
	if res := checkListItemStore.db.Find(checkListItem, chliid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCheckListItemNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return checkListItem, nil
}
