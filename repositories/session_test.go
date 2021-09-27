package repositories

import (
	"backendServer/models"
	"backendServer/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	sessionTestData, _ = utils.FillTestData(10, 5, 100)
	sessionRepo        = &SessionStore{data: sessionTestData}
)

func getSomeUser(data *models.Data) (user models.User) {
	for _, someUser := range data.Users {
		user = someUser
		return
	}
	return
}

func TestCreateSessionRepository(t *testing.T) {
	t.Parallel()

	expectedSessionRepo := &SessionStore{data: sessionTestData}

	require.Equal(t, expectedSessionRepo, CreateSessionRepository(sessionTestData))
}

func TestSessionRepositoryCreateSuccess(t *testing.T) {
	t.Parallel()

	user := getSomeUser(sessionTestData)
	SID, _ := sessionRepo.Create(user)

	require.Equal(t, user.ID, sessionTestData.Sessions[SID])
}

func TestSessionRepositoryCreateFail(t *testing.T) {
	t.Parallel()

	user := getSomeUser(sessionTestData)
	userWithWrongPassword := user
	userWithWrongPassword.Password = user.Password + "FAKE"

	_, err := sessionRepo.Create(userWithWrongPassword)

	require.Error(t, err)
}

func TestSessionRepositoryGetSuccess(t *testing.T) {
	t.Parallel()

	sessionValue := "someValue"
	user := getSomeUser(sessionTestData)
	sessionTestData.Sessions[sessionValue] = user.ID

	require.Equal(t, user, sessionRepo.Get(sessionValue))
}

func TestSessionRepositoryGetFail(t *testing.T) {
	t.Parallel()

	sessionValue := "someBadValue"

	require.Equal(t, models.User{}, sessionRepo.Get(sessionValue))
}

func TestSessionRepositoryDeleteSuccess(t *testing.T) {
	t.Parallel()

	sessionValue := "someExistingValue"
	user := getSomeUser(sessionTestData)
	sessionTestData.Sessions[sessionValue] = user.ID

	err := sessionRepo.Delete(sessionValue)
	_, notDeleted := sessionTestData.Sessions[sessionValue]
	if err != nil {
		notDeleted = true
	}

	require.Equal(t, true, !notDeleted)
}

func TestSessionRepositoryDeleteFail(t *testing.T) {
	t.Parallel()

	sessionValue := "someBadValue"

	require.Error(t, sessionRepo.Delete(sessionValue))
}
