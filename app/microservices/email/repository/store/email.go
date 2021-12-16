package store

import (
	"backendServer/app/api/models"
	"backendServer/app/microservices/email/repository"

	"gorm.io/gorm"
)

type EmailStore struct {
	db *gorm.DB
}

func CreateEmailRepository(db *gorm.DB) repository.EmailRepository {
	return &EmailStore{db: db}
}

func (emailStore *EmailStore) GetAllCards() (cards *[]models.Card, err error) {
	cards = new([]models.Card)
	err = emailStore.db.Find(cards).Error
	return
}

func (emailStore *EmailStore) GetAssignedUsers(cid uint) (users *[]models.PublicUserInfo, err error) {
	users = new([]models.PublicUserInfo)
	err = emailStore.db.Model(&models.Card{CID: cid}).Association("Users").Find(users)
	return
}

func (emailStore *EmailStore) FindBoardTitleByID(bid uint) (boardTitle string, err error) {
	board := new(models.Board)
	err = emailStore.db.First(board, bid).Error
	if err == nil {
		boardTitle = board.Title
	}
	return
}
