package repositories

import (
	"backendServer/app/api/models"
	"mime/multipart"
)

type UserRepository interface {
	Create(user *models.User) (err error)
	Update(user *models.User) (err error)
	UpdateAvatar(user *models.User, avatar *multipart.FileHeader) (err error)
	GetByLogin(login string) (user *models.User, err error)
	GetByID(uid uint) (user *models.User, err error)
	GetUserTeams(uid uint) (teams *[]models.Team, err error)
	AddUserToTeam(uid, tid uint) (err error)
	IsUserExist(user *models.User) (bool, error)
	IsEmailUsed(user *models.User) (bool, error)
	IsBoardAccessed(uid uint, bid uint) (isAccessed bool, err error)
	IsCardListAccessed(uid uint, clid uint) (isAccessed bool, err error)
	IsCardAccessed(uid uint, cid uint) (isAccessed bool, err error)
	IsCommentAccessed(uid uint, cmid uint) (isAccessed bool, err error)
	IsCheckListAccessed(uid uint, chlid uint) (isAccessed bool, err error)
	IsCheckListItemAccessed(uid uint, chliid uint) (isAccessed bool, err error)
}
