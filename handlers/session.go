package handlers

import (
	"backendServer/models"
	"backendServer/utils"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SessionHandler struct {
	SessionURL	string
	Data		*models.Data
}

func CreateSessionHandler(router *gin.RouterGroup, sessionURL string, data *models.Data) {
	handler := &SessionHandler{
		SessionURL: sessionURL,
		Data:		data,
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
		c.JSON(http.StatusBadRequest, gin.H{"error": errors.New("Bad request")})
		return
	}

	//TODO валидация данных

	sessionHandler.Data.Mu.RLock()
	user, ok := sessionHandler.Data.Users[json.Login]
	sessionHandler.Data.Mu.RUnlock()

	if !ok || user.Password != json.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Bad input data")})
		return
	}

	SID := utils.RandStringRunes(32)

	sessionHandler.Data.Mu.Lock()
	sessionHandler.Data.Sessions[SID] = user.ID
	sessionHandler.Data.Mu.Unlock()

	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   SID,
		Expires: time.Now().Add(24 * time.Hour),
		Secure: true,
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

	sessionHandler.Data.Mu.RLock()
	userID, ok := sessionHandler.Data.Sessions[session.Value]
	sessionHandler.Data.Mu.RUnlock()

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "you aren't logged in"})
		return
	}

	sessionHandler.Data.Mu.RLock()
	users := sessionHandler.Data.Users
	sessionHandler.Data.Mu.RUnlock()

	for _, user := range users {
		if user.ID == userID {
			c.IndentedJSON(http.StatusOK, user.Login)
			return
		}
	}
}

func (sessionHandler *SessionHandler) Delete(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Not authorized")})
		return
	}

	sessionHandler.Data.Mu.RLock()
	_, ok := sessionHandler.Data.Sessions[session.Value]
	sessionHandler.Data.Mu.RUnlock()

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": errors.New("Not authorized")})
		return
	}

	sessionHandler.Data.Mu.Lock()
	delete(sessionHandler.Data.Sessions, session.Value)
	sessionHandler.Data.Mu.Unlock()

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(c.Writer, session)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged out"})
}
