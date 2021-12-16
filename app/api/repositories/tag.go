package repositories

import "backendServer/app/api/models"

type TagRepository interface {
	Create(tag *models.Tag) (err error)
	Update(tag *models.Tag) (err error)
	Delete(tgid uint) (err error)
	GetByID(tgid uint) (tag *models.Tag, err error)
	AddTagToCard(uid, cid uint) (err error)
}
