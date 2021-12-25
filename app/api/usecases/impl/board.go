package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"errors"
	"time"

	"github.com/google/uuid"
)

type BoardUseCaseImpl struct {
	boardRepository     repositories.BoardRepository
	userRepository      repositories.UserRepository
	teamRepository      repositories.TeamRepository
	cardListRepository  repositories.CardListRepository
	cardRepository      repositories.CardRepository
	checkListRepository repositories.CheckListRepository
}

func CreateBoardUseCase(
	boardRepository repositories.BoardRepository,
	userRepository repositories.UserRepository,
	teamRepository repositories.TeamRepository,
	cardListRepository repositories.CardListRepository,
	cardRepository repositories.CardRepository,
	checkListRepository repositories.CheckListRepository,
) usecases.BoardUseCase {
	return &BoardUseCaseImpl{
		boardRepository:     boardRepository,
		userRepository:      userRepository,
		teamRepository:      teamRepository,
		cardListRepository:  cardListRepository,
		cardRepository:      cardRepository,
		checkListRepository: checkListRepository,
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
			return nil, err
		}
		members, err := boardUseCase.teamRepository.GetTeamMembers(team.TID)
		if err != nil {
			return nil, err
		}
		(*teams)[i].Boards = *boards
		(*teams)[i].Users = *members
	}

	toggledBoards, err := boardUseCase.userRepository.GetUserToggledBoards(uid)
	if err != nil {
		return
	}

	if len(*toggledBoards) > 0 {
		additionalTeam := models.Team{
			Title:  "Остальные доски",
			Boards: *toggledBoards,
			Type:   models.InvitedBoardsTeam,
		}

		*teams = append(*teams, additionalTeam)
	}

	return
}

func (boardUseCase *BoardUseCaseImpl) CreateBoard(board *models.Board) (bid uint, err error) {
	board.AccessPath = uuid.NewString()
	err = boardUseCase.boardRepository.Create(board)
	if err != nil {
		return 0, err
	}
	return board.BID, nil
}

func (boardUseCase *BoardUseCaseImpl) GetBoard(uid, bid uint) (board *models.Board, err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
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

	board.AvailableColors = models.AvailableColors

	members, err := boardUseCase.boardRepository.GetBoardMembers(board)
	if err != nil {
		return
	}
	board.Members = *members

	invitedMembers, err := boardUseCase.boardRepository.GetBoardInvitedMembers(bid)
	if err != nil {
		return
	}
	board.InvitedMembers = *invitedMembers

	lists, err := boardUseCase.boardRepository.GetBoardCardLists(bid)
	if err != nil {
		return nil, err
	}

	tags, err := boardUseCase.boardRepository.GetBoardTags(bid)
	if err != nil {
		return
	}
	for i, tag := range *tags {
		(*tags)[i].Color = models.AvailableColors[tag.ColorID-1]
	}
	board.Tags = *tags

	for i, list := range *lists {
		var cards *[]models.Card
		cards, err = boardUseCase.cardListRepository.GetCardListCards(list.CLID)
		if err != nil {
			return
		}
		for j, card := range *cards {
			var comments *[]models.Comment
			comments, err = boardUseCase.cardRepository.GetCardComments(card.CID)
			if err != nil {
				return
			}
			var attachments *[]models.Attachment
			attachments, err = boardUseCase.cardRepository.GetCardAttachments(card.CID)
			if err != nil {
				return
			}
			var users *[]models.PublicUserInfo
			users, err = boardUseCase.cardRepository.GetAssignedUsers(card.CID)
			if err != nil {
				return
			}

			for index, comment := range *comments {
				var user *models.PublicUserInfo
				user, err = boardUseCase.userRepository.GetPublicData(comment.UID)
				if err != nil {
					return
				}
				(*comments)[index].User = *user
				(*comments)[index].DateParsed = comment.Date.Round(time.Second).String()
			}

			(*cards)[j].Assignees = *users
			(*cards)[j].Comments = *comments
			(*cards)[j].Attachments = *attachments

			var tags *[]models.Tag
			tags, err = boardUseCase.cardRepository.GetCardTags(card.CID)
			if err != nil {
				return
			}

			for index, tag := range *tags {
				(*tags)[index].Color = models.AvailableColors[tag.ColorID-1]
			}

			(*cards)[j].Tags = *tags

			var checkLists *[]models.CheckList
			checkLists, err = boardUseCase.cardRepository.GetCardCheckLists(card.CID)
			if err != nil {
				return
			}

			for index, checkList := range *checkLists {
				var checkListItems *[]models.CheckListItem
				checkListItems, err = boardUseCase.checkListRepository.GetCheckListItems(checkList.CHLID)
				if err != nil {
					return
				}
				(*checkLists)[index].CheckListItems = *checkListItems
			}

			(*cards)[j].CheckLists = *checkLists
		}
		(*lists)[i].Cards = *cards
	}
	board.CardLists = *lists

	return
}

func (boardUseCase *BoardUseCaseImpl) UpdateBoard(uid uint, board *models.Board) (err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, board.BID)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}

	return boardUseCase.boardRepository.Update(board)
}

func (boardUseCase *BoardUseCaseImpl) DeleteBoard(uid, bid uint) (err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return err
	}
	if !isAccessed {
		return customErrors.ErrNoAccess
	}

	users, err := boardUseCase.boardRepository.GetBoardInvitedMembers(bid)
	if err != nil {
		return
	}

	for _, user := range *users {
		if isAssigned, _ := boardUseCase.userRepository.IsBoardAccessed(user.UID, bid); isAssigned {
			err = boardUseCase.userRepository.AddUserToBoard(user.UID, bid)
			if err != nil {
				return
			}
		}
	}

	return boardUseCase.boardRepository.Delete(bid)
}

func (boardUseCase *BoardUseCaseImpl) ToggleUser(uid, bid, toggledUserID uint) (board *models.Board, err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	err = boardUseCase.userRepository.AddUserToBoard(toggledUserID, bid)
	if err != nil {
		return
	}

	return boardUseCase.GetBoard(uid, bid)
}

func (boardUseCase *BoardUseCaseImpl) UpdateAccessPath(uid, bid uint) (newAccessPath string, err error) {
	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	return boardUseCase.boardRepository.UpdateAccessPath(bid)
}

func (boardUseCase *BoardUseCaseImpl) AddUserViaLink(uid uint, accessPath string) (board *models.Board, err error) {
	bid, err := boardUseCase.boardRepository.FindBoardIDByPath(accessPath)
	if err != nil {
		return
	}

	isAccessed, err := boardUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil && !errors.Is(err, customErrors.ErrNoAccess) {
		return
	}
	if !isAccessed {
		err = boardUseCase.userRepository.AddUserToBoard(uid, bid)
		if err != nil {
			return
		}
	}

	return boardUseCase.GetBoard(uid, bid)
}
