package handlers

import (
	"backendServer/app/models"
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	SessionURL     string
	SessionUseCase usecases.SessionUseCase
}

func CreateSessionHandler(router *gin.RouterGroup,
	sessionURL string,
	sessionUseCase usecases.SessionUseCase,
	mw SessionMiddleware) {
	handler := &SessionHandler{
		SessionURL:     sessionURL,
		SessionUseCase: sessionUseCase,
	}

	sessions := router.Group(handler.SessionURL)
	{
		sessions.POST("", handler.Create)
		sessions.GET("", handler.Get, mw.CheckAuth())
		sessions.DELETE("", handler.Delete, mw.CheckAuth())
	}
}

func (sessionHandler *SessionHandler) Create(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBindJSON(user); err != nil {
		_ = c.Error(errors.ErrBadRequest)
		return
	}

	SID, err := sessionHandler.SessionUseCase.Create(user)
	if err != nil {
		_ = c.Error(err)
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
	sid, exists := c.Get("sid")
	if !exists {
		_ = c.Error(errors.ErrNotAuthorized)
		return
	}

	userLogin, err := sessionHandler.SessionUseCase.Get(sid.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, userLogin)
}

func (sessionHandler *SessionHandler) Delete(c *gin.Context) {
	sid, exists := c.Get("sid")
	if !exists {
		_ = c.Error(errors.ErrNotAuthorized)
		return
	}

	err := sessionHandler.SessionUseCase.Delete(sid.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	session, _ := c.Request.Cookie("session_id")
	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(c.Writer, session)

	c.JSON(http.StatusOK, gin.H{"status": "you are logged out"})
}
