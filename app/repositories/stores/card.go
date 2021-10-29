package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"

	"gorm.io/gorm"
)

type CardStore struct {
	db *gorm.DB
}

func CreateCardRepository(db *gorm.DB) repositories.CardRepository {
	return &CardStore{db: db}
}

func (cardStore *CardStore) Create(card *models.Card) (err error) {
	// TODO
	return
}

func (cardStore *CardStore) Update(card *models.Card) (err error) {
	// TODO
	return
}

func (cardStore *CardStore) Delete(cid uint) (err error) {
	// TODO
	return
}

func (cardStore *CardStore) GetByID(cid uint) (card *models.Card, err error) {
	// TODO
	return
}
