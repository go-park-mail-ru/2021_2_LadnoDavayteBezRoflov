package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type CheckListUseCaseImpl struct {
	checkListRepository repositories.CheckListRepository
	userRepository      repositories.UserRepository
}

func CreateCheckListUseCase(
	checkListRepository repositories.CheckListRepository,
	userRepository repositories.UserRepository,
) usecases.CheckListUseCase {
	return &CheckListUseCaseImpl{
		checkListRepository: checkListRepository,
		userRepository:      userRepository,
	}
}

func (checkListUseCase *CheckListUseCaseImpl) CreateCheckList(checkList *models.CheckList) (chlid uint, err error) {
	err = checkListUseCase.checkListRepository.Create(checkList)
	if err != nil {
		return 0, err
	}
	return checkList.CHLID, nil
}

func (checkListUseCase *CheckListUseCaseImpl) GetCheckList(uid, chlid uint) (checkList *models.CheckList, err error) {
	isAccessed, err := checkListUseCase.userRepository.IsCheckListAccessed(uid, chlid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	checkList, err = checkListUseCase.checkListRepository.GetByID(chlid)
	if err != nil {
		return
	}

	checkListItems, err := checkListUseCase.checkListRepository.GetCheckListItems(chlid)
	if err != nil {
		return
	}

	checkList.CheckListItems = *checkListItems
	return
}

func (checkListUseCase *CheckListUseCaseImpl) UpdateCheckList(uid uint, checkList *models.CheckList) (err error) {
	isAccessed, err := checkListUseCase.userRepository.IsCheckListAccessed(uid, checkList.CHLID)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return checkListUseCase.checkListRepository.Update(checkList)
}

func (checkListUseCase *CheckListUseCaseImpl) DeleteCheckList(uid, chlid uint) (err error) {
	isAccessed, err := checkListUseCase.userRepository.IsCheckListAccessed(uid, chlid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return checkListUseCase.checkListRepository.Delete(chlid)
}
