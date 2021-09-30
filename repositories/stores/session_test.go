package stores

import (
	"backendServer/models"
	"backendServer/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testData, _ = utils.FillTestData(10, 5, 100)
	sessionRepo = &SessionStore{data: testData}
)

func TestCreateSessionRepository(t *testing.T) {
	t.Parallel()

	expectedSessionRepo := &SessionStore{data: testData}

	require.Equal(t, expectedSessionRepo, CreateSessionRepository(testData))
}

func TestSessionRepositoryCreateSuccess(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(testData)
	SID, _ := sessionRepo.Create(user)

	testData.Mu.RLock()
	actualSessionValue := testData.Sessions[SID]
	testData.Mu.RUnlock()

	require.Equal(t, user.ID, actualSessionValue)
}

func TestSessionRepositoryCreateFail(t *testing.T) {
	t.Parallel()

	user := utils.GetSomeUser(testData)
	userWithWrongPassword := user
	userWithWrongPassword.Password = user.Password + "FAKE"

	_, err := sessionRepo.Create(userWithWrongPassword)

	require.Error(t, err)
}

func TestSessionRepositoryGetSuccess(t *testing.T) {
	t.Parallel()

	sessionValue := "someValue"
	user := utils.GetSomeUser(testData)

	testData.Mu.Lock()
	testData.Sessions[sessionValue] = user.ID
	testData.Mu.Unlock()

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
	user := utils.GetSomeUser(testData)

	testData.Mu.Lock()
	testData.Sessions[sessionValue] = user.ID
	testData.Mu.Unlock()

	err := sessionRepo.Delete(sessionValue)
	testData.Mu.RLock()
	_, notDeleted := testData.Sessions[sessionValue]
	testData.Mu.RUnlock()
	if err != nil {
		notDeleted = true
	}

	require.True(t, !notDeleted)
}

func TestSessionRepositoryDeleteFail(t *testing.T) {
	t.Parallel()

	sessionValue := "someBadValue"

	require.Error(t, sessionRepo.Delete(sessionValue))
}
