package handlers

import (
	"backendServer/repositories"
	"errors"
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
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Not authorized")})
		return
	}

	user := boardHandler.SessionRepository.Get(session.Value)
	if user.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Not authorized")})
		return
	}

	teams := boardHandler.BoardRepository.GetAll(user.Teams)

	c.IndentedJSON(http.StatusOK, teams)
}
