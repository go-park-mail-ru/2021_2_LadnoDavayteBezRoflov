package repositories

import (
	"backendServer/models"
	"backendServer/utils"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

var (
	userTestData, _ = utils.FillTestData(10, 5, 100)
	userRepo        = &UserStore{data: userTestData}
)

func TestCreateUserRepository(t *testing.T) {
	t.Parallel()

	expectedUserRepo := &UserStore{data: userTestData}

	require.Equal(t, expectedUserRepo, CreateUserRepository(userTestData))
}

func TestUserRepositoryCreateSuccess(t *testing.T) {
	t.Parallel()

	newUser := &models.User{}
	faker.FakeData(newUser)

	user, err := userRepo.Create(*newUser)
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, userTestData.Users[newUser.Login], user)
}

func TestUserRepositoryCreateFail(t *testing.T) {
	t.Parallel()

	existingUser := getSomeUser(userTestData)

	_, errUserIsExist := userRepo.Create(existingUser)

	newUser := &models.User{}
	faker.FakeData(newUser)
	newUser.Email = existingUser.Email
	_, errEmailIsUsed := userRepo.Create(*newUser)

	require.Error(t, errUserIsExist)
	require.Error(t, errEmailIsUsed)
}
