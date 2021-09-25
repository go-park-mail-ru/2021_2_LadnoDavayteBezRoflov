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
	users := boardHandler.Data.Users
	allTeams := boardHandler.Data.Teams
	boardHandler.Data.Mu.RUnlock()

	var teamsID []uint

	for _, user := range users {
		if user.ID == userID {
			teamsID = user.Teams
			break
		}
	}

	var teams []models.Team
	for _, teamID := range teamsID {
		teams = append(teams, allTeams[teamID])
	}

	c.IndentedJSON(http.StatusOK, teams)
}
