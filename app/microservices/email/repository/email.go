package repository

import "backendServer/app/api/models"

type EmailRepository interface {
	GetAllCards() (cards *[]models.Card, err error)
	FindBoardTitleByID(bid uint) (boardTitle string, err error)
	GetAssignedUsers(cid uint) (users *[]models.PublicUserInfo, err error)
}
