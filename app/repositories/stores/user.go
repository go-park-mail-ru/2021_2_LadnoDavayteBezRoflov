package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"errors"

	_ "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserStore{db: db}
}

func (userStore *UserStore) Create(user *models.User) (err error) {
	isUserExist, err := userStore.IsUserExist(user)
	if !isUserExist {
		return
	}

	isEmailUsed, err := userStore.IsEmailUsed(user)
	if !isEmailUsed {
		return
	}

	user.HashedPassword, err = hasher.HashPassword(user.Password)
	if err != nil {
		return
	}

	err = userStore.db.Create(user).Error
	return
}

func (userStore *UserStore) Update(user *models.User) (err error) {
	oldUser, err := userStore.GetByID(user.UID)
	if err != nil {
		return
	}

	if user.Login != "" && user.Login != oldUser.Login {
		var isNewLoginExist bool
		isNewLoginExist, err = userStore.IsUserExist(user)
		if !isNewLoginExist {
			return
		}
		oldUser.Login = user.Login
	}

	if user.Email != "" && user.Email != oldUser.Email {
		var isNewEmailUsed bool
		isNewEmailUsed, err = userStore.IsEmailUsed(user)
		if !isNewEmailUsed {
			return
		}
		oldUser.Email = user.Email
	}

	if user.Password != "" && hasher.IsPasswordsEqual(user.Password, oldUser.HashedPassword) {
		oldUser.HashedPassword, err = hasher.HashPassword(user.Password)
	}

	if user.Description != "" && user.Description != oldUser.Description {
		oldUser.Description = user.Description
	}

	// TODO добавление аватара

	return userStore.db.Save(oldUser).Error
}

func (userStore *UserStore) GetByLogin(login string) (*models.User, error) {
	user := new(models.User)
	err := userStore.db.Where("login = ?", login).First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = customErrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	user := new(models.User)
	err := userStore.db.First(user, uid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = customErrors.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (userStore *UserStore) GetUserTeams(uid uint) (teams *[]models.Team, err error) {
	teams = new([]models.Team)
	err = userStore.db.Model(&models.User{UID: uid}).Association("Teams").Find(teams)

	return
}

func (userStore *UserStore) AddUserToTeam(uid, tid uint) (err error) {
	return userStore.db.Model(&models.Team{TID: tid}).Association("Users").Append(userStore.GetByID(uid))
}

func (userStore *UserStore) IsUserExist(user *models.User) (bool, error) {
	if err := userStore.db.Where("login = ?", user.Login).Find(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = customErrors.ErrUserAlreadyCreated
		}
		return true, err
	}
	return false, nil
}

func (userStore *UserStore) IsEmailUsed(user *models.User) (bool, error) {
	if err := userStore.db.Where("email = ?", user.Email).Find(user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = customErrors.ErrEmailAlreadyUsed
		}
		return true, err
	}
	return false, nil
}
