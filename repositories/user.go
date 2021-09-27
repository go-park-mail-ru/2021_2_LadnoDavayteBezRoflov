package repositories

import "backendServer/models"

type UserRepository interface {
	Create(user models.User) (finalUser models.User, err error)
}
