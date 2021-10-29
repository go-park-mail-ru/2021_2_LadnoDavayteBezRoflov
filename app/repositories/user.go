package repositories

import (
	"backendServer/app/models"
)

type UserRepository interface {
	Create(user *models.User) (err error)
	Update(user *models.User) (err error)
	GetByLogin(user *models.User) (err error)
	GetByID(uid uint) (user *models.User, err error)
	GetUserTeams(uid uint) (teams *[]models.Team, err error)
	AddUserToTeam(uid, tid uint) (err error)
}
