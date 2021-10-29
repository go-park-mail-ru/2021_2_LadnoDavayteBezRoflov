package repositories

import (
	"backendServer/app/models"
)

type BoardRepository interface {
	Create(board *models.Board) (err error)
	Update(board *models.Board) (err error)
	Delete(bid uint) (err error)
	GetByID(bid uint) (board *models.Board, err error)
	GetBoardCardLists(bid uint) (cardLists *[]models.CardList, err error)
	GetBoardCards(bid uint) (cards *[]models.Card, err error)
}
