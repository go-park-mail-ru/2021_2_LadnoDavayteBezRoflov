package repositories

import "backendServer/models"

type BoardRepository interface {
	GetAll(teamsIDs []uint) (teams []models.Team)
}
