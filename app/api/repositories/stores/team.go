package stores

import (
	"backendServer/app/api/models"
	"backendServer/app/api/repositories"
	customErrors "backendServer/pkg/errors"
	"errors"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type TeamStore struct {
	db *gorm.DB
}

func CreateTeamRepository(db *gorm.DB) repositories.TeamRepository {
	return &TeamStore{db: db}
}

func (teamStore *TeamStore) Create(team *models.Team) (err error) {
	isExist, err := teamStore.IsTeamExist(team)
	if isExist {
		return
	}
	return teamStore.db.Create(team).Error
}

func (teamStore *TeamStore) Update(team *models.Team) (err error) {
	oldTeam, err := teamStore.GetByID(team.TID)
	if err != nil {
		return
	}

	if team.Title != "" && team.Title != oldTeam.Title {
		var isNewTitleExist bool
		emptyTeam := new(models.Team)
		emptyTeam.Title = team.Title
		isNewTitleExist, err = teamStore.IsTeamExist(emptyTeam)
		if isNewTitleExist {
			return
		}
		oldTeam.Title = team.Title
	}

	return teamStore.db.Save(oldTeam).Error
}

func (teamStore *TeamStore) Delete(tid uint) (err error) {
	return teamStore.db.Delete(&models.Team{}, tid).Error
}

func (teamStore *TeamStore) GetByID(tid uint) (*models.Team, error) {
	team := new(models.Team)
	err := teamStore.db.First(team, tid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = customErrors.ErrTeamNotFound
	}
	if err != nil {
		return nil, err
	}
	return team, nil
}

func (teamStore *TeamStore) GetTeamMembers(tid uint) (members *[]models.User, err error) {
	members = new([]models.User)
	err = teamStore.db.Model(&models.Team{TID: tid}).Association("Users").Find(members)
	return
}

func (teamStore *TeamStore) GetTeamBoards(tid uint) (boards *[]models.Board, err error) {
	boards = new([]models.Board)
	err = teamStore.db.Model(&models.Team{TID: tid}).Association("Boards").Find(boards)
	return
}

func (teamStore *TeamStore) IsTeamExist(team *models.Team) (bool, error) {
	result := teamStore.db.Where("title = ?", team.Title).Find(team)
	if result.Error != nil {
		return true, result.Error
	}

	if result.RowsAffected > 0 {
		return true, nil
	}
	return false, nil
}
