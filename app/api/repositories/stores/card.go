package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	"github.com/google/uuid"

	_ "gorm.io/driver/postgres" // Register postgres driver
	"gorm.io/gorm"
)

type CardStore struct {
	db *gorm.DB
}

func CreateCardRepository(db *gorm.DB) repositories.CardRepository {
	return &CardStore{db: db}
}

func (cardStore *CardStore) Create(card *models.Card) (err error) {
	var currentMaxPos int64
	err = cardStore.db.Model(&models.Card{}).Where("cl_id = ?", card.CLID).Count(&currentMaxPos).Error
	if err != nil {
		return
	}
	card.PositionOnCardList = uint(currentMaxPos) + 1
	return cardStore.db.Create(card).Error
}

func (cardStore *CardStore) Update(card *models.Card) (err error) {
	oldCard, err := cardStore.GetByID(card.CID)
	if err != nil {
		return
	}

	if card.Title != "" && card.Title != oldCard.Title {
		oldCard.Title = card.Title
	}

	if card.Description != oldCard.Description {
		oldCard.Description = card.Description
	}

	if card.PositionOnCardList != 0 && card.CLID != 0 && (card.PositionOnCardList != oldCard.PositionOnCardList || card.CLID != oldCard.CLID) {
		err = cardStore.Move(oldCard.PositionOnCardList, card.PositionOnCardList, oldCard.CLID, card.CLID)
		if err != nil {
			return
		}
		oldCard.PositionOnCardList = card.PositionOnCardList
		oldCard.CLID = card.CLID
	}

	if card.Deadline != oldCard.Deadline {
		oldCard.Deadline = card.Deadline
	}

	if card.DeadlineChecked != oldCard.DeadlineChecked {
		oldCard.DeadlineChecked = card.DeadlineChecked
	}

	return cardStore.db.Save(oldCard).Error
}

func (cardStore *CardStore) Delete(cid uint) (err error) {
	card, err := cardStore.GetByID(cid)
	if err != nil {
		return
	}
	err = cardStore.Move(card.PositionOnCardList, (^uint(0)-1)/2, card.CLID, card.CLID)
	if err != nil {
		return
	}
	err = cardStore.db.Delete(&models.Card{}, cid).Error
	return
}

func (cardStore *CardStore) GetByID(cid uint) (*models.Card, error) {
	card := new(models.Card)
	if res := cardStore.db.Find(card, cid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCardNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return card, nil
}

func (cardStore *CardStore) GetAssignedUsers(cid uint) (users *[]models.PublicUserInfo, err error) {
	users = new([]models.PublicUserInfo)
	err = cardStore.db.Model(&models.Card{CID: cid}).Association("Users").Find(users)
	return
}

func (cardStore *CardStore) GetCardComments(cid uint) (comments *[]models.Comment, err error) {
	comments = new([]models.Comment)
	err = cardStore.db.Where("c_id = ?", cid).Order("date").Find(comments).Error
	return
}

func (cardStore *CardStore) GetCardTags(cid uint) (tags *[]models.Tag, err error) {
	tags = new([]models.Tag)
	err = cardStore.db.Model(&models.Card{CID: cid}).Association("Tags").Find(tags)
	return
}

func (cardStore *CardStore) GetCardCheckLists(cid uint) (checkLists *[]models.CheckList, err error) {
	checkLists = new([]models.CheckList)
	err = cardStore.db.Where("c_id = ?", cid).Order("chl_id").Find(checkLists).Error
	return
}

func (cardStore *CardStore) GetCardAttachments(cid uint) (attachments *[]models.Attachment, err error) {
	attachments = new([]models.Attachment)
	err = cardStore.db.Where("c_id = ?", cid).Order("at_id").Find(attachments).Error
	return
}

func (cardStore *CardStore) move(from, to, clid uint, isToLeftMove bool) (err error) {
	subQuery := cardStore.db.Model(&models.Card{}).Where(
		"cl_id = ? AND position_on_card_list BETWEEN ? AND ?",
		clid,
		from,
		to,
	)

	if isToLeftMove {
		err = subQuery.UpdateColumn("position_on_card_list", gorm.Expr("position_on_card_list - ?", 1)).Error
	} else {
		err = subQuery.UpdateColumn("position_on_card_list", gorm.Expr("position_on_card_list + ?", 1)).Error
	}

	return
}

func (cardStore *CardStore) Move(fromPos, toPos, fromCardListID, toCardListID uint) (err error) {
	if fromCardListID == toCardListID {
		if fromPos > toPos {
			err = cardStore.move(toPos, fromPos-1, fromCardListID, false)
		} else {
			err = cardStore.move(fromPos+1, toPos, fromCardListID, true)
		}
	} else {
		err = cardStore.move(fromPos, (^uint(0)-1)/2 /*максимально возможное значение позиции*/, fromCardListID, true)
		if err != nil {
			return
		}
		err = cardStore.move(toPos, (^uint(0)-1)/2 /*максимально возможное значение позиции*/, toCardListID, false)
	}
	return
}

func (cardStore *CardStore) UpdateAccessPath(cid uint) (newAccessPath string, err error) {
	newAccessPath = uuid.NewString()

	oldCard, err := cardStore.GetByID(cid)
	if err != nil {
		return
	}
	oldCard.AccessPath = newAccessPath

	err = cardStore.db.Save(oldCard).Error
	return
}

func (cardStore *CardStore) FindCardByPath(accessPath string) (card *models.Card, err error) {
	card = new(models.Card)
	err = cardStore.db.Where("access_path = ?", accessPath).Take(card).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}

	if card.CID == 0 {
		err = customErrors.ErrCardNotFound
	}

	return
}
