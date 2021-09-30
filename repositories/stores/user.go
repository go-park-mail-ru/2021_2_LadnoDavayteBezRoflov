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
	userStore.data.Mu.Lock()
	defer userStore.data.Mu.Unlock()

	_, userAlreadyCreated := userStore.data.Users[user.Login]

	if userAlreadyCreated {
		err = errors.ErrUserAlreadyCreated
		return
	}

	for _, curUser := range userStore.data.Users {
		if curUser.Email == user.Email {
			err = errors.ErrEmailAlreadyUsed
			return
		}
	}

	finalUser = user
	finalUser.ID = uint(len(userStore.data.Users))
	finalUser.Teams = []uint{}

	userStore.data.Users[user.Login] = finalUser

	return
}
