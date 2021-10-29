package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"

	"gorm.io/gorm"
)

type BoardStore struct {
	db *gorm.DB
}

func CreateBoardRepository(db *gorm.DB) repositories.BoardRepository {
	return &BoardStore{db: db}
}

func (boardStore *BoardStore) Create(board *models.Board) (err error) {
	// TODO
	return
}

func (boardStore *BoardStore) Update(board *models.Board) (err error) {
	// TODO
	return
}

func (boardStore *BoardStore) Delete(bid uint) (err error) {
	// TODO
	return
}

func (boardStore *BoardStore) GetByID(bid uint) (board *models.Board, err error) {
	// TODO
	return
}

func (boardStore *BoardStore) GetBoardCardLists(bid uint) (cardLists *[]models.CardList, err error) {
	// TODO
	return
}

func (boardStore *BoardStore) GetBoardCards(bid uint) (cards *[]models.Card, err error) {
	// TODO
	return
}
