package handlers

import (
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AttachmentHandler struct {
	AttachmentURL     string
	AttachmentUseCase usecases.AttachmentUseCase
}

func CreateAttachmentHandler(router *gin.RouterGroup,
	attachmentURL string,
	attachmentUseCase usecases.AttachmentUseCase,
	mw SessionMiddleware) {
	handler := &AttachmentHandler{
		AttachmentURL:     attachmentURL,
		AttachmentUseCase: attachmentUseCase,
	}

	attachments := router.Group(handler.AttachmentURL)
	{
		attachments.GET("/:atid", mw.CheckAuth(), mw.CSRF(), handler.GetAttachment)
		attachments.POST("/:cid", mw.CheckAuth(), mw.CSRF(), handler.CreateAttachment)
		attachments.DELETE("/:atid", mw.CheckAuth(), mw.CSRF(), handler.DeleteAttachment)
	}
}

func (attachmentHandler *AttachmentHandler) GetAttachment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	atid64 := c.Param("atid")
	atid, err := strconv.ParseUint(atid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	attachment, err := attachmentHandler.AttachmentUseCase.GetAttachment(uint(atid), uid.(uint))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, attachment)
}

func (attachmentHandler *AttachmentHandler) CreateAttachment(c *gin.Context) {
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

	file, err := c.FormFile("attachment")
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	attachment, err := attachmentHandler.AttachmentUseCase.CreateAttachment(file, uint(cid), uid.(uint))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, attachment)
}

func (attachmentHandler *AttachmentHandler) DeleteAttachment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	atid64 := c.Param("atid")
	atid, err := strconv.ParseUint(atid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = attachmentHandler.AttachmentUseCase.DeleteAttachment(uint(atid), uid.(uint))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "attachment was successfully deleted"})
}
