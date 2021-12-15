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

type CheckListItemHandler struct {
	CheckListItemURL     string
	CheckListItemUseCase usecases.CheckListItemUseCase
}

func CreateCheckListItemHandler(router *gin.RouterGroup,
	checkListItemURL string,
	checkListItemUseCase usecases.CheckListItemUseCase,
	mw SessionMiddleware) {
	handler := &CheckListItemHandler{
		CheckListItemURL:     checkListItemURL,
		CheckListItemUseCase: checkListItemUseCase,
	}

	checkListItems := router.Group(handler.CheckListItemURL)
	{
		checkListItems.GET("/:chliid", mw.CheckAuth(), mw.CSRF(), handler.GetCheckListItem)
		checkListItems.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateCheckListItem)
		checkListItems.PUT("/:chliid", mw.CheckAuth(), mw.CSRF(), handler.UpdateCheckListItem)
		checkListItems.DELETE("/:chliid", mw.CheckAuth(), mw.CSRF(), handler.DeleteCheckListItem)
	}
}

func (checkListItemHandler *CheckListItemHandler) GetCheckListItem(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chliid64 := c.Param("chliid")
	chliid, err := strconv.ParseUint(chliid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	checkListItem, err := checkListItemHandler.CheckListItemUseCase.GetCheckListItem(uid.(uint), uint(chliid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	checkListItemJSON, err := checkListItem.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", checkListItemJSON)
}

func (checkListItemHandler *CheckListItemHandler) CreateCheckListItem(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	checkListItem := new(models.CheckListItem)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, checkListItem); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	chliid, err := checkListItemHandler.CheckListItemUseCase.CreateCheckListItem(checkListItem)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"chliid": chliid})
}

func (checkListItemHandler *CheckListItemHandler) UpdateCheckListItem(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chliid64 := c.Param("chliid")
	chliid, err := strconv.ParseUint(chliid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	checkListItem := new(models.CheckListItem)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, checkListItem); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	checkListItem.CHLIID = uint(chliid)

	err = checkListItemHandler.CheckListItemUseCase.UpdateCheckListItem(uid.(uint), checkListItem)
	if err != nil {
		_ = c.Error(err)
		return
	}

	checkListItemJSON, err := checkListItem.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", checkListItemJSON)
}

func (checkListItemHandler *CheckListItemHandler) DeleteCheckListItem(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chliid64 := c.Param("chliid")
	chliid, err := strconv.ParseUint(chliid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = checkListItemHandler.CheckListItemUseCase.DeleteCheckListItem(uid.(uint), uint(chliid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "check list item was successfully deleted"})
}
