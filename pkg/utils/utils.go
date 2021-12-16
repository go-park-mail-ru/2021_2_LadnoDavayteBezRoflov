package utils

import (
	"backendServer/app/api/models"
	"regexp"
)

func ValidateUserData(user *models.User, isValidationEmailNeeded bool) (isValid bool) {
	isValid = true
	regLatinSymbols := regexp.MustCompile(".*[a-zA-Z].*")

	userLoginLen := len(user.Login)
	if userLoginLen < 3 || userLoginLen > 20 || !regLatinSymbols.MatchString(user.Login) {
		isValid = false
		return
	}

	userPasswordLen := len(user.Password)
	if userPasswordLen < 6 || userPasswordLen > 25 || !regLatinSymbols.MatchString(user.Password) {
		isValid = false
		return
	}

	if isValidationEmailNeeded && !regexp.MustCompile(".+@.+").MatchString(user.Email) {
		isValid = false
		return
	}

	return
}
