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

type CardListHandler struct {
	CardListURL     string
	CardListUseCase usecases.CardListUseCase
}

func CreateCardListHandler(router *gin.RouterGroup,
	cardListURL string,
	cardListUseCase usecases.CardListUseCase,
	mw SessionMiddleware) {
	handler := &CardListHandler{
		CardListURL:     cardListURL,
		CardListUseCase: cardListUseCase,
	}

	cards := router.Group(handler.CardListURL)
	{
		cards.GET("/:clid", mw.CheckAuth(), mw.CSRF(), handler.GetCardList)
		cards.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateCardList)
		cards.PUT("/:clid", mw.CheckAuth(), mw.CSRF(), handler.UpdateCardList)
		cards.DELETE("/:clid", mw.CheckAuth(), mw.CSRF(), handler.DeleteCardList)
	}
}

func (cardListHandler *CardListHandler) GetCardList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	clid64 := c.Param("clid")
	clid, err := strconv.ParseUint(clid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	cardList, err := cardListHandler.CardListUseCase.GetCardList(uid.(uint), uint(clid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	cardListJSON, err := cardList.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", cardListJSON)
}

func (cardListHandler *CardListHandler) CreateCardList(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	cardList := new(models.CardList)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, cardList); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	clid, err := cardListHandler.CardListUseCase.CreateCardList(cardList)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"clid": clid})
}

func (cardListHandler *CardListHandler) UpdateCardList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	clid64 := c.Param("clid")
	clid, err := strconv.ParseUint(clid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	cardList := new(models.CardList)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, cardList); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	cardList.CLID = uint(clid)
	err = cardListHandler.CardListUseCase.UpdateCardList(uid.(uint), cardList)
	if err != nil {
		_ = c.Error(err)
		return
	}

	cardListJSON, err := cardList.MarshalJSON()
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.Data(http.StatusOK, "application/json; charset=utf-8", cardListJSON)
}

func (cardListHandler *CardListHandler) DeleteCardList(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	clid64 := c.Param("clid")
	clid, err := strconv.ParseUint(clid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	err = cardListHandler.CardListUseCase.DeleteCardList(uid.(uint), uint(clid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "card list was successfully deleted"})
}
