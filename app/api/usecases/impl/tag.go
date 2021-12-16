package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type TagUseCaseImpl struct {
	tagRepository  repositories.TagRepository
	userRepository repositories.UserRepository
}

func CreateTagUseCase(
	tagRepository repositories.TagRepository,
	userRepository repositories.UserRepository,
) usecases.TagUseCase {
	return &TagUseCaseImpl{
		tagRepository:  tagRepository,
		userRepository: userRepository,
	}
}

func (tagUseCase *TagUseCaseImpl) CreateTag(tag *models.Tag) (tgid uint, err error) {
	err = tagUseCase.tagRepository.Create(tag)
	if err != nil {
		return
	}
	return tag.TGID, nil
}

func (tagUseCase *TagUseCaseImpl) GetTag(uid, tgid uint) (tag *models.Tag, err error) {
	tag, err = tagUseCase.tagRepository.GetByID(tgid)
	if err != nil {
		return
	}

	isAccessed, err := tagUseCase.userRepository.IsBoardAccessed(uid, tag.BID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return
}

func (tagUseCase *TagUseCaseImpl) UpdateTag(uid uint, tag *models.Tag) (err error) {
	isAccessed, err := tagUseCase.userRepository.IsBoardAccessed(uid, tag.BID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return tagUseCase.tagRepository.Update(tag)
}

func (tagUseCase *TagUseCaseImpl) DeleteTag(uid, tgid uint) (err error) {
	tag, err := tagUseCase.tagRepository.GetByID(tgid)
	if err != nil {
		return
	}

	isAccessed, err := tagUseCase.userRepository.IsBoardAccessed(uid, tag.BID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return tagUseCase.tagRepository.Delete(tgid)
}
