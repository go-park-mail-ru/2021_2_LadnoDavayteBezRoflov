package repositories

import (
	"backendServer/app/models"
)

type BoardRepository interface {
	GetAll(uid uint) (teams *[]models.Team, err error)
}
