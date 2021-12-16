package repositories

import (
	"backendServer/app/api/models"
)

type BoardRepository interface {
	Create(board *models.Board) (err error)
	Update(board *models.Board) (err error)
	Delete(bid uint) (err error)
	GetByID(bid uint) (board *models.Board, err error)
	GetBoardMembers(board *models.Board) (members *[]models.PublicUserInfo, err error)
	GetBoardInvitedMembers(bid uint) (members *[]models.PublicUserInfo, err error)
	GetBoardCardLists(bid uint) (cardLists *[]models.CardList, err error)
	GetBoardTags(bid uint) (tags *[]models.Tag, err error)
	GetBoardCards(bid uint) (cards *[]models.Card, err error)
	UpdateAccessPath(bid uint) (newAccessPath string, err error)
	FindBoardIDByPath(accessPath string) (bid uint, err error)
}
