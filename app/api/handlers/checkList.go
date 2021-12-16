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

type CheckListHandler struct {
	CheckListURL     string
	CheckListUseCase usecases.CheckListUseCase
}

func CreateCheckListHandler(router *gin.RouterGroup,
	checkListURL string,
	checkListUseCase usecases.CheckListUseCase,
	mw SessionMiddleware) {
	handler := &CheckListHandler{
		CheckListURL:     checkListURL,
		CheckListUseCase: checkListUseCase,
	}

	checkLists := router.Group(handler.CheckListURL)
	{
		checkLists.GET("/:chlid", mw.CheckAuth(), mw.CSRF(), handler.GetCheckList)
		checkLists.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateCheckList)
		checkLists.PUT("/:chlid", mw.CheckAuth(), mw.CSRF(), handler.UpdateCheckList)
		checkLists.DELETE("/:chlid", mw.CheckAuth(), mw.CSRF(), handler.DeleteCheckList)
	}
}

func (checkListHandler *CheckListHandler) GetCheckList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chlid64 := c.Param("chlid")
	chlid, err := strconv.ParseUint(chlid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	checkList, err := checkListHandler.CheckListUseCase.GetCheckList(uid.(uint), uint(chlid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	checkListJSON, err := checkList.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", checkListJSON)
}

func (checkListHandler *CheckListHandler) CreateCheckList(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	checkList := new(models.CheckList)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, checkList); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	chlid, err := checkListHandler.CheckListUseCase.CreateCheckList(checkList)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"chlid": chlid})
}

func (checkListHandler *CheckListHandler) UpdateCheckList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chlid64 := c.Param("chlid")
	chlid, err := strconv.ParseUint(chlid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	checkList := new(models.CheckList)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, checkList); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	checkList.CHLID = uint(chlid)

	err = checkListHandler.CheckListUseCase.UpdateCheckList(uid.(uint), checkList)
	if err != nil {
		_ = c.Error(err)
		return
	}

	checkListJSON, err := checkList.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", checkListJSON)
}

func (checkListHandler *CheckListHandler) DeleteCheckList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	chlid64 := c.Param("chlid")
	chlid, err := strconv.ParseUint(chlid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = checkListHandler.CheckListUseCase.DeleteCheckList(uid.(uint), uint(chlid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "check list was successfully deleted"})
}
