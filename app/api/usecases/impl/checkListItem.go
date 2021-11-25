package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type CheckListItemUseCaseImpl struct {
	checkListItemRepository repositories.CheckListItemRepository
	userRepository          repositories.UserRepository
}

func CreateCheckListItemUseCase(
	checkListItemRepository repositories.CheckListItemRepository,
	userRepository repositories.UserRepository,
) usecases.CheckListItemUseCase {
	return &CheckListItemUseCaseImpl{
		checkListItemRepository: checkListItemRepository,
		userRepository:          userRepository,
	}
}

func (checkListItemUseCase *CheckListItemUseCaseImpl) CreateCheckListItem(checkListItem *models.CheckListItem) (chliid uint, err error) {
	err = checkListItemUseCase.checkListItemRepository.Create(checkListItem)
	if err != nil {
		return 0, err
	}
	return checkListItem.CHLIID, nil
}

func (checkListItemUseCase *CheckListItemUseCaseImpl) GetCheckListItem(uid, chliid uint) (checkListItem *models.CheckListItem, err error) {
	isAccessed, err := checkListItemUseCase.userRepository.IsCheckListItemAccessed(uid, chliid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return checkListItemUseCase.checkListItemRepository.GetByID(chliid)
}

func (checkListItemUseCase *CheckListItemUseCaseImpl) UpdateCheckListItem(uid uint, checkListItem *models.CheckListItem) (err error) {
	isAccessed, err := checkListItemUseCase.userRepository.IsCheckListItemAccessed(uid, checkListItem.CHLIID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return checkListItemUseCase.checkListItemRepository.Update(checkListItem)
}

func (checkListItemUseCase *CheckListItemUseCaseImpl) DeleteCheckListItem(uid, chliid uint) (err error) {
	isAccessed, err := checkListItemUseCase.userRepository.IsCheckListItemAccessed(uid, chliid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return checkListItemUseCase.checkListItemRepository.Delete(chliid)
}
