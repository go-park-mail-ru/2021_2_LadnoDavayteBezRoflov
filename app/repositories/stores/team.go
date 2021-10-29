package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"

	"gorm.io/gorm"
)

type TeamStore struct {
	db *gorm.DB
}

func CreateTeamRepository(db *gorm.DB) repositories.TeamRepository {
	return &TeamStore{db: db}
}

func (teamStore *TeamStore) Create(team *models.Team) (err error) {
	// TODO
	return
}

func (teamStore *TeamStore) Update(team *models.Team) (err error) {
	// TODO
	return
}

func (teamStore *TeamStore) Delete(tid uint) (err error) {
	// TODO
	return
}

func (teamStore *TeamStore) GetByID(tid uint) (team *models.Team, err error) {
	// TODO
	return
}

func (teamStore *TeamStore) GetTeamMembers(tid uint) (members *[]models.User, err error) {
	// TODO
	return
}

func (teamStore *TeamStore) GetTeamBoards(tid uint) (boards *[]models.Board, err error) {
	boards = new([]models.Board)
	err = teamStore.db.Model(&models.Team{TID: tid}).Association("Boards").Find(boards)
	return
}
