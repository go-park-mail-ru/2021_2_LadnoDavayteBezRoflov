package usecases

import (
	"backendServer/app/models"
)

type BoardUseCase interface {
	GetAll(uid uint) (*[]models.Team, error)
}
