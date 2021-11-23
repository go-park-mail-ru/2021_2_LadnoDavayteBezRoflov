package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	"backendServer/pkg/errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	BoardURL     string
	BoardUseCase usecases.BoardUseCase
}

func CreateBoardHandler(router *gin.RouterGroup,
	boardURL string,
	boardUseCase usecases.BoardUseCase,
	mw SessionMiddleware) {
	handler := &BoardHandler{
		BoardURL:     boardURL,
		BoardUseCase: boardUseCase,
	}

	boards := router.Group(handler.BoardURL)
	{
		boards.GET("", mw.CheckAuth(), mw.CSRF(), handler.GetAllUserBoards)
		boards.POST("", mw.CheckAuth(), mw.CSRF(), handler.CreateBoard)
		boards.GET("/:bid", mw.CheckAuth(), mw.CSRF(), handler.GetBoard)
		boards.PUT("/:bid", mw.CheckAuth(), mw.CSRF(), handler.UpdateBoard)
		boards.DELETE("/:bid", mw.CheckAuth(), mw.CSRF(), handler.DeleteBoard)
		boards.PUT("/:bid/toggleuser/:uid", mw.CheckAuth(), mw.CSRF(), handler.ToggleUser)
	}
}

func (boardHandler *BoardHandler) GetAllUserBoards(c *gin.Context) {
	uid, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	boards, err := boardHandler.BoardUseCase.GetUserBoards(uid.(uint))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, boards)
}

func (boardHandler *BoardHandler) GetBoard(c *gin.Context) {
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

	board, err := boardHandler.BoardUseCase.GetBoard(uid.(uint), uint(bid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, board)
}

func (boardHandler *BoardHandler) CreateBoard(c *gin.Context) {
	_, exists := c.Get("uid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	board := new(models.Board)
	if err := c.ShouldBindJSON(board); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	bid, err := boardHandler.BoardUseCase.CreateBoard(board)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"bid": bid})
}

func (boardHandler *BoardHandler) UpdateBoard(c *gin.Context) {
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

	board := new(models.Board)
	if err := c.ShouldBindJSON(board); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}
	board.BID = uint(bid)

	err = boardHandler.BoardUseCase.UpdateBoard(uid.(uint), board)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, board)
}

func (boardHandler *BoardHandler) DeleteBoard(c *gin.Context) {
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

	err = boardHandler.BoardUseCase.DeleteBoard(uid.(uint), uint(bid))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "board was successfully deleted"})
}

func (boardHandler *BoardHandler) ToggleUser(c *gin.Context) {
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

	uid64 := c.Param("uid")
	toggledUserID, err := strconv.ParseUint(uid64, 10, 32)
	if err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	board, err := boardHandler.BoardUseCase.ToggleUser(uid.(uint), uint(bid), uint(toggledUserID))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, board)
}
