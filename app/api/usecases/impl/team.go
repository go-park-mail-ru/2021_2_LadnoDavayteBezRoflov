package impl

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
)

type TeamUseCaseImpl struct {
	teamRepository repositories.TeamRepository
	userRepository repositories.UserRepository
}

func CreateTeamUseCase(teamRepository repositories.TeamRepository, userRepository repositories.UserRepository) usecases.TeamUseCase {
	return &TeamUseCaseImpl{teamRepository: teamRepository, userRepository: userRepository}
}

func (teamUseCase *TeamUseCaseImpl) CreateTeam(uid uint, team *models.Team) (tid uint, err error) {
	err = teamUseCase.teamRepository.Create(team)
	if err != nil {
		return
	}

	err = teamUseCase.userRepository.AddUserToTeam(uid, team.TID)
	if err != nil {
		return
	}

	return team.TID, nil
}

func (teamUseCase *TeamUseCaseImpl) GetTeam(uid, tid uint) (team *models.Team, err error) {
	isMember, err := teamUseCase.userRepository.IsUserInTeam(uid, tid)
	if err != nil {
		return
	}
	if !isMember {
		err = customErrors.ErrNoAccess
		return
	}

	team, err = teamUseCase.teamRepository.GetByID(tid)
	if err != nil {
		return
	}

	boards, err := teamUseCase.teamRepository.GetTeamBoards(tid)
	if err != nil {
		return
	}
	team.Boards = *boards

	members, err := teamUseCase.teamRepository.GetTeamMembers(tid)
	if err != nil {
		return
	}
	team.Users = *members

	return
}

func (teamUseCase *TeamUseCaseImpl) UpdateTeam(uid uint, team *models.Team) (err error) {
	isMember, err := teamUseCase.userRepository.IsUserInTeam(uid, team.TID)
	if err != nil {
		return
	}
	if !isMember {
		err = customErrors.ErrNoAccess
		return
	}

	return teamUseCase.teamRepository.Update(team)
}

func (teamUseCase *TeamUseCaseImpl) DeleteTeam(uid, tid uint) (err error) {
	isMember, err := teamUseCase.userRepository.IsUserInTeam(uid, tid)
	if err != nil {
		return
	}
	if !isMember {
		err = customErrors.ErrNoAccess
		return
	}

	return teamUseCase.teamRepository.Delete(tid)
}

func (teamUseCase *TeamUseCaseImpl) ToggleUser(uid, tid, toggledUserID uint) (team *models.Team, err error) {
	isMember, err := teamUseCase.userRepository.IsUserInTeam(uid, tid)
	if err != nil {
		return
	}
	if !isMember {
		err = customErrors.ErrNoAccess
		return
	}

	err = teamUseCase.userRepository.AddUserToTeam(toggledUserID, tid)
	if err != nil {
		return
	}

	team, err = teamUseCase.GetTeam(uid, tid)
	if err == customErrors.ErrNoAccess {
		return nil, nil
	}
	return
}
