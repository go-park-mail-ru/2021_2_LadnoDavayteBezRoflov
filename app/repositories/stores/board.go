package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	customErrors "backendServer/pkg/errors"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type BoardStore struct {
	db *gorm.DB
}

func CreateBoardRepository(db *gorm.DB) repositories.BoardRepository {
	return &BoardStore{db: db}
}

func (boardStore *BoardStore) Create(board *models.Board) (err error) {
	return boardStore.db.Create(board).Error
}

func (boardStore *BoardStore) Update(board *models.Board) (err error) {
	oldBoard, err := boardStore.GetByID(board.BID)
	if err != nil {
		return
	}

	if board.Title != "" && board.Title != oldBoard.Title {
		oldBoard.Title = board.Title
	}

	if board.Description != "" && board.Description != oldBoard.Description {
		oldBoard.Description = board.Description
	}

	return boardStore.db.Save(oldBoard).Error
}

func (boardStore *BoardStore) Delete(bid uint) (err error) {
	return boardStore.db.Delete(&models.Board{}, bid).Error
}

func (boardStore *BoardStore) GetByID(bid uint) (*models.Board, error) {
	board := new(models.Board)
	if res := boardStore.db.Find(board, bid); res.RowsAffected == 0 {
		err := customErrors.ErrBoardNotFound
		return nil, err
	} else if res.Error != nil {
		return nil, res.Error
	}
	return board, nil
}

func (boardStore *BoardStore) GetBoardMembers(board *models.Board) (members *[]models.PublicUserInfo, err error) {
	members = new([]models.PublicUserInfo)
	err = boardStore.db.Model(&models.Board{BID: board.BID}).Association("Users").Find(members)
	if err != nil {
		return
	}
	teamMembers := new([]models.PublicUserInfo)
	err = boardStore.db.Model(&models.Team{TID: board.TID}).Association("Users").Find(teamMembers)
	if err != nil {
		return
	}
	*members = append(*members, *teamMembers...)
	return
}

func (boardStore *BoardStore) GetBoardInvitedMembers(bid uint) (members *[]models.PublicUserInfo, err error) {
	members = new([]models.PublicUserInfo)
	err = boardStore.db.Model(&models.Board{BID: bid}).Association("Users").Find(members)
	return
}

func (boardStore *BoardStore) GetBoardCardLists(bid uint) (cardLists *[]models.CardList, err error) {
	cardLists = new([]models.CardList)
	err = boardStore.db.Where("b_id = ?", bid).Order("position_on_board").Find(cardLists).Error
	return
}

func (boardStore *BoardStore) GetBoardCards(bid uint) (cards *[]models.Card, err error) {
	cards = new([]models.Card)
	err = boardStore.db.Model(&models.Board{BID: bid}).Association("Cards").Find(cards)
	return
}
