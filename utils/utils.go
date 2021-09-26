package utils

import (
	"backendServer/models"
	"fmt"
	"sync"

	"github.com/bxcodec/faker/v3"
)

func FillTestData(teamsAmount, boardsPerTeamAmount, usersAmount int) (data *models.Data, err error) {
	data = &models.Data{
		Sessions: map[string]uint{},
		Users:    map[string]models.User{},
		Teams:    map[uint]models.Team{},
		Mu:       &sync.RWMutex{},
	}
	for i := 0; i < teamsAmount; i++ {
		team := models.Team{}
		err = faker.FakeData(&team)
		if err != nil {
			return
		}
		team.ID = uint(i)

		for j := 0; j < boardsPerTeamAmount; j++ {
			board := models.Board{}
			err = faker.FakeData(&board)
			if err != nil {
				return
			}
			board.ID = uint(j)
			team.Boards = append(team.Boards, board)
		}

		data.Teams[team.ID] = team
	}

	for i := 0; i < usersAmount; i++ {
		user := models.User{}
		err = faker.FakeData(&user)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		user.ID = uint(i)
		user.Teams = []uint{uint(i % teamsAmount), uint(i%teamsAmount + 1), uint(i%teamsAmount + 2)}
		data.Users[user.Login] = user
	}

	return
}
