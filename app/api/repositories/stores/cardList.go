package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CardListStore struct {
	db *gorm.DB
}

func CreateCardListRepository(db *gorm.DB) repositories.CardListRepository {
	return &CardListStore{db: db}
}

func (cardListStore *CardListStore) Create(cardList *models.CardList) (err error) {
	var currentMaxPos int64
	err = cardListStore.db.Model(&models.CardList{}).Where("b_id = ?", cardList.BID).Count(&currentMaxPos).Error
	if err != nil {
		return
	}
	cardList.PositionOnBoard = uint(currentMaxPos) + 1
	return cardListStore.db.Create(cardList).Error
}

func (cardListStore *CardListStore) Update(cardList *models.CardList) (err error) {
	oldCardList, err := cardListStore.GetByID(cardList.CID)
	if err != nil {
		return
	}

	if cardList.Title != "" && cardList.Title != oldCardList.Title {
		oldCardList.Title = cardList.Title
	}

	if cardList.PositionOnBoard != 0 && cardList.PositionOnBoard != oldCardList.PositionOnBoard {
		err = cardListStore.Move(oldCardList.PositionOnBoard, cardList.PositionOnBoard, oldCardList.BID)
		if err != nil {
			return
		}
		oldCardList.PositionOnBoard = cardList.PositionOnBoard
	}

	return cardListStore.db.Save(oldCardList).Error
}

func (cardListStore *CardListStore) Delete(clid uint) (err error) {
	cardList, err := cardListStore.GetByID(clid)
	if err != nil {
		return
	}
	err = cardListStore.Move(cardList.PositionOnBoard, (^uint(0)-1)/2, cardList.BID)
	if err != nil {
		return
	}
	err = cardListStore.db.Delete(&models.CardList{}, clid).Error
	return
}

func (cardListStore *CardListStore) GetByID(clid uint) (*models.CardList, error) {
	cardList := new(models.CardList)
	if res := cardListStore.db.First(cardList, clid); res.RowsAffected == 0 {
		return nil, customErrors.ErrCardListNotFound
	} else if res.Error != nil {
		return nil, res.Error
	}
	return cardList, nil
}

func (cardListStore *CardListStore) GetCardListCards(clid uint) (cards *[]models.Card, err error) {
	cards = new([]models.Card)
	err = cardListStore.db.Where("cl_id = ?", clid).Order("position_on_card_list").Find(cards).Error
	return
}

func (cardListStore *CardListStore) move(from, to, bid uint, isToLeftMove bool) (err error) {
	subQuery := cardListStore.db.Model(&models.CardList{}).Where("b_id = ? AND position_on_board BETWEEN ? AND ?",
		bid,
		from,
		to,
	)

	if isToLeftMove {
		err = subQuery.UpdateColumn("position_on_board", gorm.Expr("position_on_board - ?", 1)).Error
	} else {
		err = subQuery.UpdateColumn("position_on_board", gorm.Expr("position_on_board + ?", 1)).Error
	}

	return
}

func (cardListStore *CardListStore) Move(fromPos, toPos, bid uint) (err error) {
	if fromPos > toPos {
		err = cardListStore.move(toPos, fromPos-1, bid, false)
	} else {
		err = cardListStore.move(fromPos+1, toPos, bid, true)
	}
	return
}
