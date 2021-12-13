package usecase

import (
	"backendServer/app/api/models"

	"gopkg.in/gomail.v2"
)

type EmailUseCase interface {
	SendFirstLetter(userInfo *models.PublicUserInfo) (emailLetter *gomail.Message)
	SendNotifications() (emailLetters *[]gomail.Message, err error)
}
