package handlers

import (
	"backendServer/errors"
	"backendServer/models"
	"backendServer/repositories"
	"backendServer/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	SessionURL        string
	SessionRepository repositories.SessionRepository
}

func CreateSessionHandler(router *gin.RouterGroup, sessionURL string, sessionRepository repositories.SessionRepository) {
	handler := &SessionHandler{
		SessionURL:        sessionURL,
		SessionRepository: sessionRepository,
	}

	sessions := router.Group(handler.SessionURL)
	{
		sessions.POST("", handler.Create)
		sessions.GET("", handler.Get)
		sessions.DELETE("", handler.Delete)
	}
}

func (sessionHandler *SessionHandler) Create(c *gin.Context) {
	var json models.User
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	json.Email = "email@template"
	if !utils.ValidateUserData(json) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrBadInputData.Error()})
		return
	}

	SID, err := sessionHandler.SessionRepository.Create(json)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func (sessionHandler *SessionHandler) Get(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you aren't logged in"})
		return
	}

	user := sessionHandler.SessionRepository.Get(session.Value)

	if user.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you aren't logged in"})
		return
	}

	c.IndentedJSON(http.StatusOK, user.Login)
	return
}

func (sessionHandler *SessionHandler) Delete(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.ErrNotAuthorized.Error()})
		return
	}

	err = sessionHandler.SessionRepository.Delete(session.Value)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(c.Writer, session)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged out"})
}
