package stores

import (
	"backendServer/app/models"
	"backendServer/app/repositories"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/hasher"
	"errors"

	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func CreateUserRepository(db *gorm.DB) repositories.UserRepository {
	return &UserStore{db: db}
}

func (userStore *UserStore) Create(user *models.User) (err error) {
	if isUserExist := userStore.db.Where("login = ?", user.Login).Find(user).RowsAffected; isUserExist > 0 {
		err = customErrors.ErrUserAlreadyCreated
		return
	}

	if isEmailUsed := userStore.db.Where("email = ?", user.Email).Find(user).RowsAffected; isEmailUsed > 0 {
		err = customErrors.ErrEmailAlreadyUsed
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
	// TODO
	return
}

func (userStore *UserStore) GetByLogin(user *models.User) (err error) {
	return userStore.db.Where("login = ?", user.Login).First(user).Error
}

func (userStore *UserStore) GetByID(uid uint) (user *models.User, err error) {
	user = new(models.User)
	if err = userStore.db.First(user, uid).Error; errors.Is(err, gorm.ErrRecordNotFound) {
		err = customErrors.ErrUserNotFound
	}
	return
}

func (userStore *UserStore) GetUserTeams(uid uint) (teams *[]models.Team, err error) {
	teams = new([]models.Team)
	err = userStore.db.Model(&models.User{UID: uid}).Association("Teams").Find(teams)

	return
}

func (userStore *UserStore) AddUserToTeam(uid, tid uint) (err error) {
	return userStore.db.Model(&models.Team{TID: tid}).Association("Users").Append(userStore.GetByID(uid))
}
