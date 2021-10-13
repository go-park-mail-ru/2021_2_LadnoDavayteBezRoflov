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
		c.JSON(errors.ResolveErrorToCode(errors.ErrBadRequest), gin.H{"error": errors.ErrBadRequest.Error()})
		return
	}

	if !utils.ValidateUserData(json, false) {
		c.JSON(errors.ResolveErrorToCode(errors.ErrBadInputData), gin.H{"error": errors.ErrBadInputData.Error()})
		return
	}

	SID, err := sessionHandler.SessionRepository.Create(json)
	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    SID,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(c.Writer, cookie)
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func (sessionHandler *SessionHandler) Get(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(errors.ResolveErrorToCode(errors.ErrNotAuthorized), gin.H{"status": "you aren't logged in"})
		return
	}

	user, err := sessionHandler.SessionRepository.Get(session.Value)
	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user.Login)
	return
}

func (sessionHandler *SessionHandler) Delete(c *gin.Context) {
	session, err := c.Request.Cookie("session_id")
	if err == http.ErrNoCookie {
		c.JSON(errors.ResolveErrorToCode(errors.ErrNotAuthorized), gin.H{"error": errors.ErrNotAuthorized.Error()})
		return
	}

	err = sessionHandler.SessionRepository.Delete(session.Value)
	if err != errors.ErrInternal {
		session.Expires = time.Now().AddDate(0, 0, -1)
		http.SetCookie(c.Writer, session)
	}

	if err != nil {
		c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "you are logged out"})
}
