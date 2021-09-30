package stores

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/utils"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/stretchr/testify/require"
)

var userRepo = &UserStore{data: testData}

func TestCreateUserRepository(t *testing.T) {
	t.Parallel()

	expectedUserRepo := &UserStore{data: testData}

	require.Equal(t, expectedUserRepo, CreateUserRepository(testData))
}

func TestUserRepositoryCreateSuccess(t *testing.T) {
	t.Parallel()

	newUser := &models.User{}
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}

	user, errCreate := userRepo.Create(*newUser)
	if errCreate != nil {
		require.NoError(t, err)
	}

	testData.Mu.RLock()
	expectedUser := testData.Users[user.Login]
	testData.Mu.RUnlock()

	require.Equal(t, expectedUser, user)
}

func TestUserRepositoryCreateFail(t *testing.T) {
	t.Parallel()

	existingUser := utils.GetSomeUser(testData)

	_, errUserIsExist := userRepo.Create(existingUser)

	newUser := &models.User{}
	err := faker.FakeData(newUser)
	if err != nil {
		t.Error(err)
	}
	newUser.Email = existingUser.Email
	_, errEmailIsUsed := userRepo.Create(*newUser)

	require.Error(t, errors.ErrUserAlreadyCreated, errUserIsExist)
	require.Error(t, errors.ErrEmailAlreadyUsed, errEmailIsUsed)
}
