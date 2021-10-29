package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
)

type BoardUseCaseImpl struct {
	boardRepository repositories.BoardRepository
	userRepository  repositories.UserRepository
	teamRepository  repositories.TeamRepository
}

func CreateBoardUseCase(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
) usecases.BoardUseCase {
	return &BoardUseCaseImpl{
		boardRepository: boardRepository,
		userRepository:  userRepository,
		teamRepository:  teamRepository,
	}
}

func (boardUseCase *BoardUseCaseImpl) GetUserBoards(uid uint) (teams *[]models.Team, err error) {
	teams, err = boardUseCase.userRepository.GetUserTeams(uid)
	if err != nil {
		return
	}

	for i, team := range *teams {
		boards, boardsErr := boardUseCase.teamRepository.GetTeamBoards(team.TID)
		if boardsErr != nil {
			err = boardsErr
			return
		}
		(*teams)[i].Boards = *boards
	}

	return
}
