package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"

	"gorm.io/gorm"
)

type CardListStore struct {
	db *gorm.DB
}

func CreateCardListRepository(db *gorm.DB) repositories.CardListRepository {
	return &CardListStore{db: db}
}

func (cardListStore *CardListStore) Create(cardList *models.CardList) (err error) {
	// TODO
	return
}

func (cardListStore *CardListStore) Update(cardList *models.CardList) (err error) {
	// TODO
	return
}

func (cardListStore *CardListStore) Delete(clid uint) (err error) {
	// TODO
	return
}

func (cardListStore *CardListStore) GetByID(clid uint) (cardList *models.CardList, err error) {
	// TODO
	return
}

func (cardListStore *CardListStore) GetCardListCards(clid uint) (cards *[]models.Card, err error) {
	// TODO
	return
}
