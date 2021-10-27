package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
)

type BoardUseCaseImpl struct {
	boardRepository repositories.BoardRepository
}

func CreateBoardUseCase(boardRepository repositories.BoardRepository) usecases.BoardUseCase {
	return &BoardUseCaseImpl{boardRepository: boardRepository}
}

func (boardUseCaseImpl *BoardUseCaseImpl) GetAll(uid uint) (*[]models.Team, error) {
	teams, err := boardUseCaseImpl.boardRepository.GetAll(uid)
	if err != nil {
		return nil, err
	}
	return teams, nil
}
