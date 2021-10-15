package handlers

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserURL     string
	UserUseCase usecases.UserUseCase
}

func CreateUserHandler(router *gin.RouterGroup, userURL string, userUseCase usecases.UserUseCase) {
	handler := &UserHandler{
		UserURL:     userURL,
		UserUseCase: userUseCase,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("", handler.Create)
	}
}

func (userHandler *UserHandler) Create(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBindJSON(user); err != nil {
		c.JSON(errors.ResolveErrorToCode(errors.ErrBadRequest), gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	SID, err := userHandler.UserUseCase.Create(user)
	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusCreated, gin.H{"status": "you are logged in"})
}
