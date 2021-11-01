package impl

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	"backendServer/app/usecases"
	customErrors "backendServer/pkg/errors"
)

type BoardUseCaseImpl struct {
	boardRepository    repositories.BoardRepository
	userRepository     repositories.UserRepository
	teamRepository     repositories.TeamRepository
	cardListRepository repositories.CardListRepository
}

func CreateBoardUseCase(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
	cardListRepository repositories.CardListRepository,
) usecases.BoardUseCase {
	return &BoardUseCaseImpl{
		boardRepository:    boardRepository,
		userRepository:     userRepository,
		teamRepository:     teamRepository,
		cardListRepository: cardListRepository,
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

func (boardUseCase *BoardUseCaseImpl) CreateBoard(board *models.Board) (bid uint, err error) {
	err = boardUseCase.boardRepository.Create(board)
	if err != nil {
		return 0, err
	}
	return board.BID, nil
}

func (boardUseCase *BoardUseCaseImpl) GetBoard(uid, bid uint) (board *models.Board, err error) {
	isAccessed, err := boardUseCase.isUserHaveAccessToBoard(uid, bid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}
	board, err = boardUseCase.boardRepository.GetByID(bid)
	if err != nil {
		return
	}

	lists, err := boardUseCase.boardRepository.GetBoardCardLists(bid)
	if err != nil {
		return nil, err
	}

	for i, list := range *lists {
		cards := new([]models.Card)
		cards, err = boardUseCase.cardListRepository.GetCardListCards(list.CLID)
		if err != nil {
			return
		}
		(*lists)[i].Cards = *cards
	}
	board.CardLists = *lists

	return
}

func (boardUseCase *BoardUseCaseImpl) UpdateBoard(uid uint, board *models.Board) (err error) {
	isAccessed, err := boardUseCase.isUserHaveAccessToBoard(uid, board.BID)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}

	err = boardUseCase.boardRepository.Update(board)
	return
}

func (boardUseCase *BoardUseCaseImpl) DeleteBoard(uid, bid uint) (err error) {
	isAccessed, err := boardUseCase.isUserHaveAccessToBoard(uid, bid)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}
	return boardUseCase.boardRepository.Delete(bid)
}

func (boardUseCase *BoardUseCaseImpl) isUserHaveAccessToBoard(uid, bid uint) (isAccessed bool, err error) {
	teams, err := boardUseCase.userRepository.GetUserTeams(uid)
	if err != nil {
		return
	}

	for _, team := range *teams {
		boards, boardsErr := boardUseCase.teamRepository.GetTeamBoards(team.TID)
		if boardsErr != nil {
			err = boardsErr
			return
		}
		for _, board := range *boards {
			if board.BID == bid {
				isAccessed = true
				return
			}
		}
	}
	isAccessed = false
	return
}
