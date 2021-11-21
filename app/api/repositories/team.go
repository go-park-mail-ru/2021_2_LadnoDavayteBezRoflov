package repositories

import (
	models2 "backendServer/app/api/models"
)

type TeamRepository interface {
	Create(team *models2.Team) (err error)
	Update(team *models2.Team) (err error)
	Delete(tid uint) (err error)
	GetByID(tid uint) (team *models2.Team, err error)
	GetTeamMembers(tid uint) (members *[]models2.User, err error)
	GetTeamBoards(tid uint) (boards *[]models2.Board, err error)
	IsTeamExist(team *models2.Team) (bool, error)
}
