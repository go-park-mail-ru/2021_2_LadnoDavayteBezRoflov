package stores

import (
	"backendServer/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	boardTestData, _ = utils.FillTestData(10, 5, 100)
	boardRepo        = &BoardStore{data: userTestData}
)

func TestCreateBoardRepository(t *testing.T) {
	t.Parallel()

	expectedBoardRepo := &BoardStore{data: boardTestData}

	require.Equal(t, expectedBoardRepo, CreateBoardRepository(boardTestData))
}

func TestBoardRepository_GetAll(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(boardTestData)
	teamsIDs := user.Teams
	teams := boardRepo.GetAll(teamsIDs)

	allTeamsReceived := true
	for index, team := range teams {
		if team.ID != teamsIDs[index] {
			allTeamsReceived = false
		}
	}
	require.Equal(t, true, allTeamsReceived)
}
