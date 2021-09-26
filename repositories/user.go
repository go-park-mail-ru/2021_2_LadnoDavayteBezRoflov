package repositories

import (
	"backendServer/models"
	"errors"
)

type UserRepository struct {
	data *models.Data
}

func CreateUserRepository(data *models.Data) (userRepository UserRepository) {
	return UserRepository{data: data}
}

func (userRepository *UserRepository) Create(user models.User) (finalUser models.User, err error) {
	userRepository.data.Mu.RLock()
	_, userAlreadyCreated := userRepository.data.Users[user.Login]
	users := userRepository.data.Users
	userRepository.data.Mu.RUnlock()

	if userAlreadyCreated {
		err = errors.New("User already created")
		return
	}

	for _, curUser := range users {
		if curUser.Email == user.Email {
			err = errors.New("Email already used")
			return
		}
	}

	finalUser = user
	finalUser.ID = uint(len(users))

	userRepository.data.Mu.Lock()
	userRepository.data.Users[user.Login] = finalUser
	userRepository.data.Mu.Unlock()

	return
}
