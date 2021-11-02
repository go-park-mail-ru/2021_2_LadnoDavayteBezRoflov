package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
)

type CardUseCaseImpl struct {
	cardRepository repositories.CardRepository
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
	// TODO добавить проверку на то, что пользователь имеет доступ к этой карте

	card, err = cardUseCase.cardRepository.GetByID(cid)
	return
}

func (cardUseCase *CardUseCaseImpl) UpdateCard(uid uint, card *models.Card) (err error) {
	// TODO добавить проверку на то, что пользователь имеет доступ к этой карте

	return cardUseCase.cardRepository.Update(card)
}

func (cardUseCase *CardUseCaseImpl) DeleteCard(uid, cid uint) (err error) {
	// TODO добавить проверку на то, что пользователь имеет доступ к этой карте

	return cardUseCase.cardRepository.Delete(cid)
}
