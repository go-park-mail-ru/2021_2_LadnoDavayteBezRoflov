package stores

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
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}

	user, err := userRepo.Create(*newUser)
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, userTestData.Users[newUser.Login], user)
}

func TestUserRepositoryCreateFail(t *testing.T) {
	t.Parallel()

	existingUser := utils.GetSomeUser(userTestData)

	_, errUserIsExist := userRepo.Create(existingUser)

	newUser := &models.User{}
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}
	newUser.Email = existingUser.Email
	_, errEmailIsUsed := userRepo.Create(*newUser)

	require.Error(t, errUserIsExist)
	require.Error(t, errEmailIsUsed)
}
