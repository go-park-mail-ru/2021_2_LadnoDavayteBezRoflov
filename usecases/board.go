package usecases

import "backendServer/models"

type BoardUseCase interface {
	GetAll(uid uint) (*[]models.Team, error)
}
