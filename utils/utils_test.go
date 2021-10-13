package utils

import (
	"backendServer/models"
	"reflect"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

var (
	normalUser = models.User{
		Login:    "latinChars",
		Password: "latinCharsAgain",
		Email:    "default@email.ru",
	}
	successTestsUserData = []struct {
		testName string
		user     models.User
	}{
		{
			testName: "normal user",
		},
		{
			testName: "email without .<smth>",
			user: models.User{
				Email: "default@mail",
			},
		},
	}

	failTestsUserData = []struct {
		testName string
		user     models.User
	}{
		{
			testName: "login with digits",
			user: models.User{
				Login: "123456",
			},
		},
		{
			testName: "login with russian chars",
			user: models.User{
				Login: "РусскиеСимволы",
			},
		},
		{
			testName: "too short login",
			user: models.User{
				Login: "sh",
			},
		},
		{
			testName: "too long login",
			user: models.User{
				Login: "tooLongLoginForUserOnThisServer",
			},
		},
		{
			testName: "password with digits",
			user: models.User{
				Password: "123456",
			},
		},
		{
			testName: "password with russian chars",
			user: models.User{
				Password: "парольНаРусском",
			},
		},
		{
			testName: "too short password",
			user: models.User{
				Password: "short",
			},
		},
		{
			testName: "too long password",
			user: models.User{
				Password: "veryVeryLongPasswordForUser",
			},
		},
		{
			testName: "email without @",
			user: models.User{
				Email: "emailWithoutA",
			},
		},
	}
)

func TestFillTestData(t *testing.T) {
	t.Parallel()

	dataFilled := true
	data, _ := FillTestData(10, 10, 10)

TeamLoop:
	for _, team := range data.Teams {
		if team.Boards == nil || len(team.Boards) == 0 || team.Title == "" {
			dataFilled = false
			break
		}
		for _, board := range team.Boards {
			if board.Title == "" || board.Description == "" || board.Tasks == nil || len(board.Tasks) == 0 {
				dataFilled = false
				break TeamLoop
			}
		}
	}

	for _, user := range data.Users {
		if user.Login == "" || user.Password == "" || user.Email == "" || user.Teams == nil {
			dataFilled = false
			break
		}
	}

	require.Equal(t, true, dataFilled)
}

func TestValidateUserDataSuccess(t *testing.T) {
	t.Parallel()

	for _, tt := range successTestsUserData {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			user := normalUser
			if tt.user.Login != "" {
				user.Login = tt.user.Login
			} else if tt.user.Password != "" {
				user.Password = tt.user.Password
			} else if tt.user.Email != "" {
				user.Email = tt.user.Email
			}
			isValid := ValidateUserData(user, true)
			require.Equal(t, true, isValid)
		})
	}
}

func TestValidateUserDataFail(t *testing.T) {
	t.Parallel()

	for _, tt := range failTestsUserData {
		tt := tt
		t.Run(tt.testName, func(t *testing.T) {
			t.Parallel()

			user := normalUser
			if tt.user.Login != "" {
				user.Login = tt.user.Login
			} else if tt.user.Password != "" {
				user.Password = tt.user.Password
			} else if tt.user.Email != "" {
				user.Email = tt.user.Email
			}
			isValid := ValidateUserData(user, true)
			require.Equal(t, false, isValid)
		})
	}
}

func TestGetSomeUser(t *testing.T) {
	t.Parallel()

	data := &models.Data{}
	err := faker.FakeData(data)
	if err != nil {
		t.Error(err)
	}

	randomUser := GetSomeUser(data)
	isExist := false
	for _, user := range data.Users {
		if reflect.DeepEqual(user, randomUser) {
			isExist = true
			break
		}
	}

	require.Equal(t, true, isExist)
}
