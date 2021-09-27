package repositories

import (
	"backendServer/models"
)

type BoardRepository interface {
	GetAll(teamsIDs []uint) (teams []models.Team)
}

type BoardStore struct {
	data *models.Data
}

func CreateBoardRepository(data *models.Data) BoardRepository {
	return &BoardStore{data: data}
}

func (boardStore *BoardStore) GetAll(teamsIDs []uint) (teams []models.Team) {
	boardStore.data.Mu.RLock()
	allTeams := boardStore.data.Teams
	boardStore.data.Mu.RUnlock()

	for _, teamID := range teamsIDs {
		teams = append(teams, allTeams[teamID])
	}

	return
}
