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

func CreateUserHandler(router *gin.RouterGroup, userURL string, userUseCase usecases.UserUseCase, mw SessionMiddleware) {
	handler := &UserHandler{
		UserURL:     userURL,
		UserUseCase: userUseCase,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("", handler.CreateUser)
		users.GET("/:login", mw.CheckAuth(), handler.GetUser)
		users.PUT("/:login", mw.CheckAuth(), handler.UpdateUser)
		users.PUT("/:login/upload", mw.CheckAuth(), handler.UpdateUserAvatar)
	}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
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

func (userHandler *UserHandler) GetUser(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	login := c.Param("login")

	user, err := userHandler.UserUseCase.Get(uid.(uint), login)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	user := new(models.User)
	if err := c.ShouldBindJSON(user); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	user.UID = uid.(uint)
	err := userHandler.UserUseCase.Update(user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user.Avatar)
}

func (userHandler *UserHandler) UpdateUserAvatar(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	user := new(models.User)
	user.UID = uid.(uint)

	file, err := c.FormFile("avatar")
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	user.Avatar = file.Filename
	user.AvatarFile = *file

	err = userHandler.UserUseCase.UpdateAvatar(user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user.Avatar)
}
