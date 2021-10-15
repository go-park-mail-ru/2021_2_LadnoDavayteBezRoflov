package handlers

import (
	"backendServer/errors"
	"backendServer/usecases"
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
	mw Middleware) {
	handler := &BoardHandler{
		BoardURL:     boardURL,
		BoardUseCase: boardUseCase,
	}

	boards := router.Group(handler.BoardURL)
	{
		boards.GET("", handler.GetAll, mw.CheckAuth())
	}
}

func (boardHandler *BoardHandler) GetAll(c *gin.Context) {
	uid, isExist := c.Get("uid")
	if !isExist {
		c.JSON(errors.ResolveErrorToCode(errors.ErrNotAuthorized), gin.H{"error": errors.ErrNotAuthorized.Error()})
		return
	}

	teams, err := boardHandler.BoardUseCase.GetAll(uid.(uint))
	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, teams)
}
