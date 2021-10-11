package handlers

import (
	"backendServer/errors"
	"backendServer/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BoardHandler struct {
	BoardURL          string
	BoardRepository   repositories.BoardRepository
	SessionRepository repositories.SessionRepository
}

func CreateBoardHandler(router *gin.RouterGroup,
	boardURL string,
	boardRepository repositories.BoardRepository,
	sessionRepository repositories.SessionRepository) {
	handler := &BoardHandler{
		BoardURL:          boardURL,
		BoardRepository:   boardRepository,
		SessionRepository: sessionRepository,
	}

	boards := router.Group(handler.BoardURL)
	{
		boards.GET("", handler.GetAll)
	}
}

func (boardHandler *BoardHandler) GetAll(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(errors.ResolveErrorToCode(errors.ErrNotAuthorized), gin.H{"error": errors.ErrNotAuthorized.Error()})
		return
	}

	user, err := boardHandler.SessionRepository.Get(session.Value)
	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	teams := boardHandler.BoardRepository.GetAll(user.Teams)

	c.JSON(http.StatusOK, teams)
}
