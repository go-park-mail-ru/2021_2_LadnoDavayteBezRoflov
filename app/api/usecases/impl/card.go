package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type CardUseCaseImpl struct {
	cardRepository repositories.CardRepository
	userRepository repositories.UserRepository
}

func CreateCardUseCase(cardRepository repositories.CardRepository, userRepository repositories.UserRepository) usecases.CardUseCase {
	return &CardUseCaseImpl{cardRepository: cardRepository, userRepository: userRepository}
}

func (cardUseCase *CardUseCaseImpl) CreateCard(card *models.Card) (cid uint, err error) {
	err = cardUseCase.cardRepository.Create(card)
	if err != nil {
		return 0, err
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

	comments, err := cardUseCase.cardRepository.GetCardComments(cid)
	if err != nil {
		return
	}

	card.Comments = *comments
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
