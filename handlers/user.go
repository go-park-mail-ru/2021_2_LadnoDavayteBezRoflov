package handlers

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories"
	"backendServer/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserURL           string
	UserRepository    repositories.UserRepository
	SessionRepository repositories.SessionRepository
}

func CreateUserHandler(router *gin.RouterGroup,
	userURL string,
	userRepository repositories.UserRepository,
	sessionRepository repositories.SessionRepository) {
	handler := &UserHandler{
		UserURL:           userURL,
		UserRepository:    userRepository,
		SessionRepository: sessionRepository,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("", handler.Create)
	}
}

func (userHandler *UserHandler) Create(c *gin.Context) {
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	if !utils.ValidateUserData(json) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrBadInputData.Error()})
		return
	}

	user, userCreateErr := userHandler.UserRepository.Create(json)
	if userCreateErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": userCreateErr.Error()})
		return
	}

	SID, sessionCreateErr := userHandler.SessionRepository.Create(user)
	if sessionCreateErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": sessionCreateErr.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusCreated, gin.H{"status": "you are logged in"})
}
