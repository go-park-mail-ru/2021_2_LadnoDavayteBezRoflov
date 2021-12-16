package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"errors"

	"github.com/google/uuid"
)

type CardUseCaseImpl struct {
	cardRepository repositories.CardRepository
	userRepository repositories.UserRepository
	tagRepository  repositories.TagRepository
}

func CreateCardUseCase(
	cardRepository repositories.CardRepository,
	userRepository repositories.UserRepository,
	tagRepository repositories.TagRepository,
) usecases.CardUseCase {
	return &CardUseCaseImpl{cardRepository: cardRepository, userRepository: userRepository, tagRepository: tagRepository}
}

func (cardUseCase *CardUseCaseImpl) CreateCard(card *models.Card) (cid uint, err error) {
	card.AccessPath = uuid.NewString()
	err = cardUseCase.cardRepository.Create(card)
	if err != nil {
		return
	}
	return card.CID, nil
}

func (cardUseCase *CardUseCaseImpl) GetCard(uid, cid uint) (card *models.Card, err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	card, err = cardUseCase.cardRepository.GetByID(cid)
	if err != nil {
		return
	}

	tags, err := cardUseCase.cardRepository.GetCardTags(cid)
	if err != nil {
		return
	}
	for i, tag := range *tags {
		(*tags)[i].Color = models.AvailableColors[tag.ColorID-1]
	}
	card.Tags = *tags

	comments, err := cardUseCase.cardRepository.GetCardComments(cid)
	if err != nil {
		return
	}
	for i, comment := range *comments {
		var user *models.PublicUserInfo
		user, err = cardUseCase.userRepository.GetPublicData(comment.UID)
		if err != nil {
			return
		}
		(*comments)[i].User = *user
	}
	card.Comments = *comments

	assignees, err := cardUseCase.cardRepository.GetAssignedUsers(cid)
	if err != nil {
		return
	}
	card.Assignees = *assignees

	return
}

func (cardUseCase *CardUseCaseImpl) UpdateCard(uid uint, card *models.Card) (err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, card.CID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return cardUseCase.cardRepository.Update(card)
}

func (cardUseCase *CardUseCaseImpl) DeleteCard(uid, cid uint) (err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return cardUseCase.cardRepository.Delete(cid)
}

func (cardUseCase *CardUseCaseImpl) ToggleUser(uid, cid, toggledUserID uint) (card *models.Card, err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	err = cardUseCase.userRepository.AddUserToCard(toggledUserID, cid)
	if err != nil {
		return
	}

	return cardUseCase.GetCard(uid, cid)
}

func (cardUseCase *CardUseCaseImpl) ToggleTag(uid, cid, toggledTagID uint) (card *models.Card, err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	err = cardUseCase.tagRepository.AddTagToCard(toggledTagID, cid)
	if err != nil {
		return
	}

	return cardUseCase.GetCard(uid, cid)
}

func (cardUseCase *CardUseCaseImpl) UpdateAccessPath(uid, cid uint) (newAccessPath string, err error) {
	isAccessed, err := cardUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return cardUseCase.cardRepository.UpdateAccessPath(cid)
}

func (cardUseCase *CardUseCaseImpl) AddUserViaLink(uid uint, accessPath string) (card *models.Card, err error) {
	card, err = cardUseCase.cardRepository.FindCardByPath(accessPath)
	if err != nil {
		return
	}

	isAccessed, err := cardUseCase.userRepository.IsBoardAccessed(uid, card.BID)
	if err != nil && !errors.Is(err, customErrors.ErrNoAccess) {
		return
	}
	if !isAccessed {
		err = cardUseCase.userRepository.AddUserToBoard(uid, card.BID)
		if err != nil {
			return
		}
	}

	isAssigned, err := cardUseCase.userRepository.IsCardAssigned(uid, card.CID)
	if err != nil && !errors.Is(err, customErrors.ErrNoAccess) {
		return
	}
	if !isAssigned {
		err = cardUseCase.userRepository.AddUserToCard(uid, card.CID)
		if err != nil {
			return
		}
	}

	return
}
