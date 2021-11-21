package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	customErrors "backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	CardURL     string
	CardUseCase usecases.CardUseCase
}

func CreateCardHandler(router *gin.RouterGroup,
	cardURL string,
	cardUseCase usecases.CardUseCase,
	mw SessionMiddleware) {
	handler := &CardHandler{
		CardURL:     cardURL,
		CardUseCase: cardUseCase,
	}

	cards := router.Group(handler.CardURL)
	{
		cards.GET("/:cid", mw.CheckAuth(), mw.CSRF(), handler.GetCard)
		cards.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateCard)
		cards.PUT("/:cid", mw.CheckAuth(), mw.CSRF(), handler.UpdateCard)
		cards.DELETE("/:cid", mw.CheckAuth(), mw.CSRF(), handler.DeleteCard)
	}
}

func (cardHandler *CardHandler) GetCard(c *gin.Context) {
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

	card, err := cardHandler.CardUseCase.GetCard(uid.(uint), uint(cid))
	if err != nil {
		return
	}

	c.JSON(http.StatusOK, card)
}

func (cardHandler *CardHandler) CreateCard(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	card := new(models.Card)
	if err := c.ShouldBindJSON(card); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	cid, err := cardHandler.CardUseCase.CreateCard(card)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"cid": cid})
}

func (cardHandler *CardHandler) UpdateCard(c *gin.Context) {
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

	card := new(models.Card)
	if err := c.ShouldBindJSON(card); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	card.CID = uint(cid)

	err = cardHandler.CardUseCase.UpdateCard(uid.(uint), card)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, card)
}

func (cardHandler *CardHandler) DeleteCard(c *gin.Context) {
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

	err = cardHandler.CardUseCase.DeleteCard(uid.(uint), uint(cid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "card was successfully deleted"})
}
