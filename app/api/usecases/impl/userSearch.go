package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type UserSearchUseCaseImpl struct {
	userRepository  repositories.UserRepository
	cardRepository  repositories.CardRepository
	teamRepository  repositories.TeamRepository
	boardRepository repositories.BoardRepository
}

func CreateUserSearchUseCase(
	userRepository repositories.UserRepository,
	cardRepository repositories.CardRepository,
	teamRepository repositories.TeamRepository,
	boardRepository repositories.BoardRepository,
) usecases.UserSearchUseCase {
	return &UserSearchUseCaseImpl{
		userRepository:  userRepository,
		cardRepository:  cardRepository,
		teamRepository:  teamRepository,
		boardRepository: boardRepository,
	}
}

func (userSearchUseCase *UserSearchUseCaseImpl) FindForCard(uid, cid uint, text string) (users *[]models.UserSearchInfo, err error) {
	isAccessed, err := userSearchUseCase.userRepository.IsCardAccessed(uid, cid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	card, err := userSearchUseCase.cardRepository.GetByID(cid)
	if err != nil {
		return
	}

	matchedUsers, err := userSearchUseCase.userRepository.FindBoardMembersByLogin(card.BID, text, 15)
	if err != nil {
		return
	}

	assignedUsers, err := userSearchUseCase.cardRepository.GetAssignedUsers(cid)
	if err != nil {
		return
	}

	users = new([]models.UserSearchInfo)
	for _, matchedUser := range *matchedUsers {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, assignedUser := range *assignedUsers {
			if assignedUser.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*users = append(*users, user)
	}

	return
}

func (userSearchUseCase *UserSearchUseCaseImpl) FindForTeam(uid, tid uint, text string) (users *[]models.UserSearchInfo, err error) {
	isMember, err := userSearchUseCase.userRepository.IsUserInTeam(uid, tid)
	if err != nil {
		return
	}
	if !isMember {
		err = customErrors.ErrNoAccess
		return
	}

	matchedUsers, err := userSearchUseCase.userRepository.FindAllByLogin(text, 15)
	if err != nil {
		return
	}

	members, err := userSearchUseCase.teamRepository.GetTeamMembers(tid)
	if err != nil {
		return
	}

	users = new([]models.UserSearchInfo)
	for _, matchedUser := range *matchedUsers {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, member := range *members {
			if member.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*users = append(*users, user)
	}

	return
}

func (userSearchUseCase *UserSearchUseCaseImpl) FindForBoard(uid, bid uint, text string) (users *[]models.UserSearchInfo, err error) {
	isAccessed, err := userSearchUseCase.userRepository.IsBoardAccessed(uid, bid)
	if err != nil {
		return
	}
	if !isAccessed {
		err = customErrors.ErrNoAccess
		return
	}

	matchedUsers, err := userSearchUseCase.userRepository.FindBoardInvitedMembersByLogin(bid, text, 15)
	if err != nil {
		return
	}

	members, err := userSearchUseCase.boardRepository.GetBoardInvitedMembers(bid)
	if err != nil {
		return
	}

	users = new([]models.UserSearchInfo)
	for _, matchedUser := range *matchedUsers {
		user := models.UserSearchInfo{
			UID:    matchedUser.UID,
			Login:  matchedUser.Login,
			Avatar: matchedUser.Avatar,
		}

		for _, member := range *members {
			if member.Login == matchedUser.Login {
				user.Added = true
				break
			}
		}

		*users = append(*users, user)
	}

	return
}
