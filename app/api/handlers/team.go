package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/mailru/easyjson"

	"github.com/gin-gonic/gin"
)

type TeamHandler struct {
	TeamURL     string
	TeamUseCase usecases.TeamUseCase
}

func CreateTeamHandler(router *gin.RouterGroup,
	teamURL string,
	teamUseCase usecases.TeamUseCase,
	mw SessionMiddleware) {
	handler := &TeamHandler{
		TeamURL:     teamURL,
		TeamUseCase: teamUseCase,
	}

	teams := router.Group(handler.TeamURL)
	{
		teams.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateTeam)
		teams.GET("/:tid", mw.CheckAuth(), mw.CSRF(), handler.GetTeam)
		teams.PUT("/:tid", mw.CheckAuth(), mw.CSRF(), handler.UpdateTeam)
		teams.DELETE("/:tid", mw.CheckAuth(), mw.CSRF(), handler.DeleteTeam)
		teams.PUT("/:tid/toggleuser/:uid", mw.CheckAuth(), mw.CSRF(), handler.ToggleUser)
	}
}

func (teamHandler *TeamHandler) CreateTeam(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	team := new(models.Team)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, team); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	tid, err := teamHandler.TeamUseCase.CreateTeam(uid.(uint), team)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tid": tid})
}

func (teamHandler *TeamHandler) GetTeam(c *gin.Context) {
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

	team, err := teamHandler.TeamUseCase.GetTeam(uid.(uint), uint(tid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	teamJSON, err := team.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", teamJSON)
}

func (teamHandler *TeamHandler) UpdateTeam(c *gin.Context) {
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

	team := new(models.Team)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, team); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	team.TID = uint(tid)

	err = teamHandler.TeamUseCase.UpdateTeam(uid.(uint), team)
	if err != nil {
		_ = c.Error(err)
		return
	}

	teamJSON, err := team.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", teamJSON)
}

func (teamHandler *TeamHandler) DeleteTeam(c *gin.Context) {
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

	err = teamHandler.TeamUseCase.DeleteTeam(uid.(uint), uint(tid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "team was successfully deleted"})
}

func (teamHandler *TeamHandler) ToggleUser(c *gin.Context) {
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

	uid64 := c.Param("uid")
	toggledUserID, err := strconv.ParseUint(uid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	team, err := teamHandler.TeamUseCase.ToggleUser(uid.(uint), uint(tid), uint(toggledUserID))
	if err != nil {
		_ = c.Error(err)
		return
	}

	teamJSON, err := team.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", teamJSON)
}
