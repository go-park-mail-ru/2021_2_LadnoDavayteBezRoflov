package handlers

import (
	"backendServer/models"
	"backendServer/repositories"
	"errors"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Bad request")})
		return
	}

	// TODO валидация данных

	user, err := userHandler.UserRepository.Create(json)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
		return
	}

	SID, err := userHandler.SessionRepository.Create(user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err})
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
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}
