package stores

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories"
)

type UserStore struct {
	data *models.Data
}

func CreateUserRepository(data *models.Data) repositories.UserRepository {
	return &UserStore{data: data}
}

func (userStore *UserStore) Create(user models.User) (finalUser models.User, err error) {
	userStore.data.Mu.RLock()
	_, userAlreadyCreated := userStore.data.Users[user.Login]
	users := userStore.data.Users
	userStore.data.Mu.RUnlock()

	if userAlreadyCreated {
		err = errors.ErrUserAlreadyCreated
		return
	}

	for _, curUser := range users {
		if curUser.Email == user.Email {
			err = errors.ErrEmailAlreadyUsed
			return
		}
	}

	finalUser = user
	finalUser.ID = uint(len(users))

	userStore.data.Mu.Lock()
	userStore.data.Users[user.Login] = finalUser
	userStore.data.Mu.Unlock()

	return
}
