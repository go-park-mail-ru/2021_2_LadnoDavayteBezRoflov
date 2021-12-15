package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mailru/easyjson"
)

type TagHandler struct {
	TagURL     string
	TagUseCase usecases.TagUseCase
}

func CreateTagHandler(router *gin.RouterGroup,
	tagURL string,
	tagUseCase usecases.TagUseCase,
	mw SessionMiddleware) {
	handler := &TagHandler{
		TagURL:     tagURL,
		TagUseCase: tagUseCase,
	}

	tags := router.Group(handler.TagURL)
	{
		tags.GET("/:tgid", mw.CheckAuth(), mw.CSRF(), handler.GetTag)
		tags.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateTag)
		tags.PUT("/:tgid", mw.CheckAuth(), mw.CSRF(), handler.UpdateTag)
		tags.DELETE("/:tgid", mw.CheckAuth(), mw.CSRF(), handler.DeleteTag)
	}
}

func (tagHandler *TagHandler) GetTag(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	tgid64 := c.Param("tgid")
	tgid, err := strconv.ParseUint(tgid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	tag, err := tagHandler.TagUseCase.GetTag(uid.(uint), uint(tgid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	tagJSON, err := tag.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", tagJSON)
}

func (tagHandler *TagHandler) CreateTag(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	tag := new(models.Tag)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, tag); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	tgid, err := tagHandler.TagUseCase.CreateTag(tag)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"tgid": tgid})
}

func (tagHandler *TagHandler) UpdateTag(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	tgid64 := c.Param("tgid")
	tgid, err := strconv.ParseUint(tgid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	tag := new(models.Tag)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, tag); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	tag.TGID = uint(tgid)

	err = tagHandler.TagUseCase.UpdateTag(uid.(uint), tag)
	if err != nil {
		_ = c.Error(err)
		return
	}

	tagJSON, err := tag.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", tagJSON)
}

func (tagHandler *TagHandler) DeleteTag(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	tgid64 := c.Param("tgid")
	tgid, err := strconv.ParseUint(tgid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = tagHandler.TagUseCase.DeleteTag(uid.(uint), uint(tgid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "tag was successfully deleted"})
}
