package handlers

import (
	"backendServer/app/api/models"
	"backendServer/app/api/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/sessionCookieController"
	"net/http"
	"time"

	"github.com/mailru/easyjson"

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
		sessions.GET("", mw.CheckAuth(), mw.CSRF(), handler.Get)
		sessions.DELETE("", mw.CheckAuth(), mw.CSRF(), handler.Delete)
	}
}

func (sessionHandler *SessionHandler) Create(c *gin.Context) {
	user := new(models.User)
	if err := easyjson.UnmarshalFromReader(c.Request.Body, user); err != nil {
		_ = c.Error(customErrors.ErrBadRequest)
		return
	}

	sid, err := sessionHandler.SessionUseCase.Create(user)
	if err != nil {
		_ = c.Error(err)
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}
	http.SetCookie(c.Writer, sessionCookieController.CreateSessionCookie(sid))
	c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
}

func (sessionHandler *SessionHandler) Get(c *gin.Context) {
	sid, exists := c.Get("sid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	userLogin, err := sessionHandler.SessionUseCase.Get(sid.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"login": userLogin})
}

func (sessionHandler *SessionHandler) Delete(c *gin.Context) {
	sid, exists := c.Get("sid")
	if !exists {
		_ = c.Error(customErrors.ErrNotAuthorized)
		return
	}

	err := sessionHandler.SessionUseCase.Delete(sid.(string))
	if err != nil {
		_ = c.Error(err)
		return
	}

	session, _ := c.Request.Cookie("session_id")
	sessionCookieController.SetSessionCookieExpired(session)
	http.SetCookie(c.Writer, session)

	csrf, _ := c.Request.Cookie("csrf_token")
	csrf.Path = "/"
	csrf.Expires = time.Now().Add(-1)
	http.SetCookie(c.Writer, csrf)

	c.JSON(http.StatusOK, gin.H{"status": "you are logged out"})
}
