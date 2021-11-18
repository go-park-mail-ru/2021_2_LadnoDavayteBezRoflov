package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	customErrors "backendServer/pkg/errors"

	"gorm.io/gorm"
)

type CheckListStore struct {
	db *gorm.DB
}

func CreateCheckListRepository(db *gorm.DB) repositories.CheckListRepository {
	return &CheckListStore{db: db}
}

func (checkListStore *CheckListStore) Create(checkList *models.CheckList) (err error) {
	return checkListStore.db.Create(checkList).Error
}

func (checkListStore *CheckListStore) Update(checkList *models.CheckList) (err error) {
	oldCheckList, err := checkListStore.GetByID(checkList.CHLID)
	if err != nil {
		return
	}

	if checkList.Title != "" && checkList.Title != oldCheckList.Title {
		oldCheckList.Title = checkList.Title
	}

	return checkListStore.db.Save(oldCheckList).Error
}

func (checkListStore *CheckListStore) Delete(chlid uint) (err error) {
	return checkListStore.db.Delete(&models.CheckList{}, chlid).Error
}

func (checkListStore *CheckListStore) GetByID(chlid uint) (*models.CheckList, error) {
	checkList := new(models.CheckList)
	if res := checkListStore.db.Find(checkList, chlid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCheckListNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return checkList, nil
}

func (checkListStore *CheckListStore) GetCheckListItems(chlid uint) (checkListItems *[]models.CheckListItem, err error) {
	checkListItems = new([]models.CheckListItem)
	err = checkListStore.db.Where("chl_id = ?", chlid).Find(checkListItems).Error
	return
}
