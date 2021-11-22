package usecases

import (
	"backendServer/app/models"
)

type UserSearchUseCase interface {
	FindForCard(uid, cid uint, text string) (users *[]models.UserSearchInfo, err error)
	FindForTeam(uid, tid uint, text string) (users *[]models.UserSearchInfo, err error)
	FindForBoard(uid, bid uint, text string) (users *[]models.UserSearchInfo, err error)
}
