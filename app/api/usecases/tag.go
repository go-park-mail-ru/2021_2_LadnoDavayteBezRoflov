package usecases

import "backendServer/app/api/models"

type TagUseCase interface {
	CreateTag(tag *models.Tag) (tgid uint, err error)
	GetTag(uid, tgid uint) (tag *models.Tag, err error)
	UpdateTag(uid uint, tag *models.Tag) (err error)
	DeleteTag(uid, tgid uint) (err error)
}
