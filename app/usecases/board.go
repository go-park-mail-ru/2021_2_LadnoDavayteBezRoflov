package usecases

import (
	"backendServer/app/models"
)

type BoardUseCase interface {
	GetUserBoards(uid uint) (*[]models.Team, error)
	CreateBoard(board *models.Board) (bid uint, err error)
	GetBoard(uid, bid uint) (board *models.Board, err error)
	UpdateBoard(uid uint, board *models.Board) (err error)
	DeleteBoard(uid, bid uint) (err error)
	ToggleUser(uid, bid, toggledUserID uint) (board *models.Board, err error)
}
