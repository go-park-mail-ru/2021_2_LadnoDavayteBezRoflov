package repositories

import (
	models2 "backendServer/app/api/models"
)

type BoardRepository interface {
	Create(board *models2.Board) (err error)
	Update(board *models2.Board) (err error)
	Delete(bid uint) (err error)
	GetByID(bid uint) (board *models2.Board, err error)
	GetBoardCardLists(bid uint) (cardLists *[]models2.CardList, err error)
	GetBoardCards(bid uint) (cards *[]models2.Card, err error)
}
