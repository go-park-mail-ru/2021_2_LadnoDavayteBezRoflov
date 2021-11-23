package repositories

import (
	"backendServer/app/api/models"
)

type TeamRepository interface {
	Create(team *models.Team) (err error)
	Update(team *models.Team) (err error)
	Delete(tid uint) (err error)
	GetByID(tid uint) (team *models.Team, err error)
	GetTeamMembers(tid uint) (members *[]models.User, err error)
	GetTeamBoards(tid uint) (boards *[]models.Board, err error)
	IsTeamExist(team *models.Team) (bool, error)
}
