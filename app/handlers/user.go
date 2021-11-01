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
	AvatarsPath string
	UserUseCase usecases.UserUseCase
}

func CreateUserHandler(router *gin.RouterGroup, userURL string, userUseCase usecases.UserUseCase, mw SessionMiddleware, avatarsPath string) {
	handler := &UserHandler{
		UserURL:     userURL,
		AvatarsPath: avatarsPath,
		UserUseCase: userUseCase,
	}

	users := router.Group(handler.UserURL)
	{
		users.POST("", mw.CheckAuth(), handler.CreateUser)
		users.GET("/:login", mw.CheckAuth(), handler.GetUser)
		users.POST("/:login", mw.CheckAuth(), handler.UpdateUser)
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

	login := c.Param("bid")

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

	login := c.Param("bid")

	user := new(models.User)
	if err := c.ShouldBindJSON(user); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	type Passwords struct {
		NewPassword string `json:"password_repeat"`
		OldPassword string `json:"old_password"`
	}
	passwords := &Passwords{}
	if err := c.ShouldBindJSON(passwords); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	if user.Avatar != "" {
		resp, err := http.Get(user.Avatar)
		if err != nil {
			_ = c.Error(customErrors.ErrBadRequest)
			return
		}
		defer func(err error) {
			if err != nil {
				_ = c.Error(err)
			}
		}(resp.Body.Close())

		user.AvatarFile = resp.Body
	}

	user.UID = uid.(uint)
	user.AvatarsPath = userHandler.AvatarsPath
	err := userHandler.UserUseCase.Update(login, passwords.NewPassword, passwords.OldPassword, user)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user.Avatar)
}
