package repositories

import "backendServer/models"

type BoardRepository interface {
	GetAll(uid uint) (teams *[]models.Team, err error)
}
