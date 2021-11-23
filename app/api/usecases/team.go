package usecases

import "backendServer/app/api/models"

type TeamUseCase interface {
	CreateTeam(uid uint, team *models.Team) (tid uint, err error)
	GetTeam(uid, tid uint) (team *models.Team, err error)
	UpdateTeam(uid uint, team *models.Team) (err error)
	DeleteTeam(uid, tid uint) (err error)
	ToggleUser(uid, tid, toggledUserID uint) (team *models.Team, err error)
}
