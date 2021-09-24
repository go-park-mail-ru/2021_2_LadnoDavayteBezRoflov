package handlers

import (
	"backendServer/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BoardHandler struct {
	BoardURL	string
	Data		*models.Data
}

func CreateBoardHandler(router *gin.RouterGroup, boardURL string, data *models.Data) {
	handler := &BoardHandler{
		BoardURL:	boardURL,
		Data:		data,
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

	boardHandler.Data.Mu.RLock()
	userID, ok := boardHandler.Data.Sessions[session.Value]
	boardHandler.Data.Mu.RUnlock()

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Not authorized")})
		return
	}

	boardHandler.Data.Mu.RLock()
	boards := boardHandler.Data.Boards[userID]
	boardHandler.Data.Mu.RUnlock()

	c.IndentedJSON(http.StatusOK, boards)
}
