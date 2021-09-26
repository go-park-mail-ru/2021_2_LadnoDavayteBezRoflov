package repositories

import (
	"backendServer/models"
)

type BoardRepository struct {
	data *models.Data
}

func CreateBoardRepository(data *models.Data) (boardRepository BoardRepository) {
	return BoardRepository{data: data}
}

func (boardRepository *BoardRepository) GetAll(teamsIDs []uint) (teams []models.Team) {
	boardRepository.data.Mu.RLock()
	allTeams := boardRepository.data.Teams
	boardRepository.data.Mu.RUnlock()

	for _, teamID := range teamsIDs {
		teams = append(teams, allTeams[teamID])
	}

	return
}
