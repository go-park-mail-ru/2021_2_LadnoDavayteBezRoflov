package stores

import (
	"backendServer/models"
	"backendServer/repositories"
)

type BoardStore struct {
	data *models.Data
}

func CreateBoardRepository(data *models.Data) repositories.BoardRepository {
	return &BoardStore{data: data}
}

func (boardStore *BoardStore) GetAll(uid uint) (teams *[]models.Team, err error) {
	boardStore.data.Mu.RLock()
	defer boardStore.data.Mu.RUnlock()

	// TODO Временно закомментировано для того, чтобы можно было просматривать доски с нового пользователя
	/*
		for _, teamID := range teamsIDs {
			teams = append(teams, allTeams[teamID])
		}*/

	for _, team := range boardStore.data.Teams {
		*teams = append(*teams, team)
	}

	return
}
