package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	customErrors "backendServer/pkg/errors"
)

type CardUseCaseImpl struct {
	cardRepository repositories.CardRepository
	userRepository repositories.UserRepository
}

func CreateCardUseCase(cardRepository repositories.CardRepository) usecases.CardUseCase {
	return &CardUseCaseImpl{cardRepository: cardRepository}
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
