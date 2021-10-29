package usecases

import (
	"backendServer/app/models"
)

type BoardUseCase interface {
	GetUserBoards(uid uint) (*[]models.Team, error)
}
