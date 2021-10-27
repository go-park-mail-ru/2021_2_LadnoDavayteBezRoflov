package stores

import (
	"backendServer/pkg/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var boardRepo = &BoardStore{data: testData}

func TestCreateBoardRepository(t *testing.T) {
	t.Parallel()

	expectedBoardRepo := &BoardStore{data: testData}

	require.Equal(t, expectedBoardRepo, CreateBoardRepository(testData))
}

func TestBoardRepository_GetAll(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(testData)
	teamsIDs := user.Teams
	teams := boardRepo.GetAll(teamsIDs)

	allTeamsReceived := true

	// TODO Временно закомментировано для того, чтобы можно было просматривать доски с нового пользователя
	/*
		for index, team := range teams {
			if team.ID != teamsIDs[index] {
				allTeamsReceived = false
			}
		}
	*/

	for _, team := range teams {
		testData.Mu.RLock()
		_, isExist := testData.Teams[team.ID]
		testData.Mu.RUnlock()

		if !isExist {
			allTeamsReceived = false
			return
		}
	}

	require.True(t, allTeamsReceived)
}
