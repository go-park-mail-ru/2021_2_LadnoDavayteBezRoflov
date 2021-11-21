package repositories

import (
	models2 "backendServer/app/api/models"
	"mime/multipart"
)

type UserRepository interface {
	Create(user *models2.User) (err error)
	Update(user *models2.User) (err error)
	UpdateAvatar(user *models2.User, avatar *multipart.FileHeader) (err error)
	GetByLogin(login string) (user *models2.User, err error)
	GetByID(uid uint) (user *models2.User, err error)
	GetUserTeams(uid uint) (teams *[]models2.Team, err error)
	AddUserToTeam(uid, tid uint) (err error)
	IsUserExist(user *models2.User) (bool, error)
	IsEmailUsed(user *models2.User) (bool, error)
	IsBoardAccessed(uid uint, bid uint) (isAccessed bool, err error)
	IsCardListAccessed(uid uint, clid uint) (isAccessed bool, err error)
	IsCardAccessed(uid uint, cid uint) (isAccessed bool, err error)
	IsCommentAccessed(uid uint, cmid uint) (isAccessed bool, err error)
	IsCheckListAccessed(uid uint, chlid uint) (isAccessed bool, err error)
	IsCheckListItemAccessed(uid uint, chliid uint) (isAccessed bool, err error)
}
