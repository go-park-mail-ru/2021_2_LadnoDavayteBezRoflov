package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"io"
	"os"
	"strings"

	"github.com/google/uuid"

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
	if isUserExist {
		return
	}

	isEmailUsed, err := userStore.IsEmailUsed(user)
	if isEmailUsed {
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

	if user.Avatar != "" {
		fileName := uuid.NewString()
		out := new(os.File)
		out, err = os.Create(strings.Join([]string{user.AvatarsPath, "/", fileName}, ""))
		defer func(closeErr error) {
			if closeErr != nil {
				err = closeErr
			}
		}(out.Close())

		_, err = io.Copy(out, user.AvatarFile)
		if err != nil {
			return
		}

		user.Avatar = fileName
		oldUser.Avatar = user.Avatar
	}

	return userStore.db.Save(oldUser).Error
}

func (userStore *UserStore) GetByLogin(login string) (*models.User, error) {
	user := new(models.User)
	if res := userStore.db.Where("login = ?", login).First(user); res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, customErrors.ErrUserNotFound
	}
	return user, nil
}

func (userStore *UserStore) GetByID(uid uint) (*models.User, error) {
	user := new(models.User)
	if res := userStore.db.First(user, uid); res.Error != nil {
		return nil, res.Error
	} else if res.RowsAffected == 0 {
		return nil, customErrors.ErrUserNotFound
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
	if res := userStore.db.Select("login").Where("login = ?", user.Login).Find(user); res.Error != nil {
		return true, res.Error
	} else if res.RowsAffected == 0 {
		return false, nil
	}
	return true, customErrors.ErrUserAlreadyCreated
}

func (userStore *UserStore) IsEmailUsed(user *models.User) (bool, error) {
	if res := userStore.db.Select("email").Where("email = ?", user.Email).Find(user); res.Error != nil {
		return true, res.Error
	} else if res.RowsAffected == 0 {
		return false, nil
	}
	return true, customErrors.ErrEmailAlreadyUsed
}
