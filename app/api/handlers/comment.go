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

type CommentHandler struct {
	CommentURL     string
	CommentUseCase usecases.CommentUseCase
}

func CreateCommentHandler(router *gin.RouterGroup,
	commentURL string,
	commentUseCase usecases.CommentUseCase,
	mw SessionMiddleware) {
	handler := &CommentHandler{
		CommentURL:     commentURL,
		CommentUseCase: commentUseCase,
	}

	comments := router.Group(handler.CommentURL)
	{
		comments.GET("/:cmid", mw.CheckAuth(), mw.CSRF(), handler.GetComment)
		comments.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateComment)
		comments.PUT("/:cmid", mw.CheckAuth(), mw.CSRF(), handler.UpdateComment)
		comments.DELETE("/:cmid", mw.CheckAuth(), mw.CSRF(), handler.DeleteComment)
	}
}

func (commentHandler *CommentHandler) GetComment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	cmid64 := c.Param("cmid")
	cmid, err := strconv.ParseUint(cmid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	comment, err := commentHandler.CommentUseCase.GetComment(uid.(uint), uint(cmid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	commentJSON, err := comment.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", commentJSON)
}

func (commentHandler *CommentHandler) CreateComment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	comment := new(models.Comment)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, comment); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	comment.UID = uid.(uint)

	comment, err := commentHandler.CommentUseCase.CreateComment(comment)
	if err != nil {
		_ = c.Error(err)
		return
	}

	commentJSON, err := comment.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", commentJSON)
}

func (commentHandler *CommentHandler) UpdateComment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	cmid64 := c.Param("cmid")
	cmid, err := strconv.ParseUint(cmid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	comment := new(models.Comment)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, comment); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	comment.CMID = uint(cmid)

	err = commentHandler.CommentUseCase.UpdateComment(uid.(uint), comment)
	if err != nil {
		_ = c.Error(err)
		return
	}

	commentJSON, err := comment.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", commentJSON)
}

func (commentHandler *CommentHandler) DeleteComment(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	cmid64 := c.Param("cmid")
	cmid, err := strconv.ParseUint(cmid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = commentHandler.CommentUseCase.DeleteComment(uid.(uint), uint(cmid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "comment was successfully deleted"})
}
