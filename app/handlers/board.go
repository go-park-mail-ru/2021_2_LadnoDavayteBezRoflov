package handlers

import (
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"net/http"

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
		boards.GET("", mw.CheckAuth(), handler.GetAllUserBoards)
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
