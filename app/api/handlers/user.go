package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/sessionCookieController"
	"net/http"

	"github.com/mailru/easyjson"

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
		users.GET("/:login", mw.CheckAuth(), mw.CSRF(), handler.GetUser)
		users.PUT("/:login", mw.CheckAuth(), mw.CSRF(), handler.UpdateUser)
		users.PUT("/:login/upload", mw.CheckAuth(), mw.CSRF(), handler.UpdateUserAvatar)
	}
}

func (userHandler *UserHandler) CreateUser(c *gin.Context) {
	user := new(models.User)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, user); err != nil {
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

	userJSON, err := user.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}

func (userHandler *UserHandler) UpdateUser(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	user := new(models.User)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, user); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	user.UID = uid.(uint)
	err := userHandler.UserUseCase.Update(user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
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

	err = userHandler.UserUseCase.UpdateAvatar(user, file)
	if err != nil {
		_ = c.Error(err)
		return
	}

	userJSON, err := user.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", userJSON)
}
