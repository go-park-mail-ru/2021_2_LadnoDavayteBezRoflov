package handlers

import (
	"backendServer/app/models"
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/sessionCookieController"
	"net/http"

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
	user := new(models.User)
	if err := c.ShouldBindJSON(user); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	sid, err := userHandler.UserUseCase.Create(user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	http.SetCookie(c.Writer, sessionCookieController.CreateSessionCookie(sid))
	c.JSON(http.StatusCreated, gin.H{"status": "you are logged in"})
}
