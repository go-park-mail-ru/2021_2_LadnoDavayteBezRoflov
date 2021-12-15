package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserSearchHandler struct {
	UserSearchURL     string
	UserSearchUseCase usecases.UserSearchUseCase
}

func CreateUserSearchHandler(router *gin.RouterGroup,
	userSearchURL string,
	userSearchUseCase usecases.UserSearchUseCase,
	mw SessionMiddleware) {
	handler := &UserSearchHandler{
		UserSearchURL:     userSearchURL,
		UserSearchUseCase: userSearchUseCase,
	}

	userSearch := router.Group(handler.UserSearchURL)
	{
		userSearch.GET("/card/:cid/:text", mw.CheckAuth(), mw.CSRF(), handler.FindForCard)
		userSearch.GET("/team/:tid/:text", mw.CheckAuth(), mw.CSRF(), handler.FindForTeam)
		userSearch.GET("/board/:bid/:text", mw.CheckAuth(), mw.CSRF(), handler.FindForBoard)
	}
}

func (userSearchHandler *UserSearchHandler) FindForCard(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	cid64 := c.Param("cid")
	cid, err := strconv.ParseUint(cid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	users, err := userSearchHandler.UserSearchUseCase.FindForCard(uid.(uint), uint(cid), c.Param("text"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	usersInfo := new(models.UsersSearchInfo)
	*usersInfo = *users

	usersInfoJSON, err := usersInfo.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", usersInfoJSON)
}

func (userSearchHandler *UserSearchHandler) FindForTeam(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	tid64 := c.Param("tid")
	tid, err := strconv.ParseUint(tid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	users, err := userSearchHandler.UserSearchUseCase.FindForTeam(uid.(uint), uint(tid), c.Param("text"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	usersInfo := new(models.UsersSearchInfo)
	*usersInfo = *users

	usersInfoJSON, err := usersInfo.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", usersInfoJSON)
}

func (userSearchHandler *UserSearchHandler) FindForBoard(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	bid64 := c.Param("bid")
	bid, err := strconv.ParseUint(bid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	users, err := userSearchHandler.UserSearchUseCase.FindForBoard(uid.(uint), uint(bid), c.Param("text"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	usersInfo := new(models.UsersSearchInfo)
	*usersInfo = *users

	usersInfoJSON, err := usersInfo.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", usersInfoJSON)
}
