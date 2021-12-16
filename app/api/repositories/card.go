package repositories

import (
	"backendServer/app/api/models"
)

type CardRepository interface {
	Create(card *models.Card) (err error)
	Update(card *models.Card) (err error)
	Delete(cid uint) (err error)
	GetByID(cid uint) (card *models.Card, err error)
	GetAssignedUsers(cid uint) (users *[]models.PublicUserInfo, err error)
	GetCardComments(cid uint) (comments *[]models.Comment, err error)
	GetCardTags(cid uint) (tags *[]models.Tag, err error)
	GetCardCheckLists(cid uint) (checkLists *[]models.CheckList, err error)
	GetCardAttachments(cid uint) (attachments *[]models.Attachment, err error)
	Move(fromPos, toPos, fromCardListID, toCardListID uint) (err error)
	UpdateAccessPath(cid uint) (newAccessPath string, err error)
	FindCardByPath(accessPath string) (card *models.Card, err error)
}
