package handlers

import (
	"backendServer/app/usecases"
	"backendServer/pkg/sessionCookieController"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SessionMiddleware interface {
	CheckAuth() gin.HandlerFunc
}

type SessionMiddlewareImpl struct {
	sessionUseCase usecases.SessionUseCase
}

func CreateSessionMiddleware(sessionUseCase usecases.SessionUseCase) SessionMiddleware {
	return &SessionMiddlewareImpl{sessionUseCase: sessionUseCase}
}

func (middleware *SessionMiddlewareImpl) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Request.Cookie("session_id")
		if err != nil {
			return
		}

		sid := session.Value
		uid, err := middleware.sessionUseCase.GetUID(sid)
		if err != nil {
			sessionCookieController.SetSessionCookieExpired(session)
			http.SetCookie(c.Writer, session)
			_ = c.Error(err)
			return
		}

		if sessionCookieController.IsSessionCookieExpiresSoon(session) {
			err := middleware.sessionUseCase.AddTime(sid, uint(sessionCookieController.SessionCookieLifeTimeInSecs.Seconds()))
			if err != nil {
				_ = c.Error(err)
				return
			}
			sessionCookieController.UpdateSessionCookieExpires(session)
			http.SetCookie(c.Writer, session)
		}

		c.Set("uid", uid)
		c.Set("sid", sid)

		return
	}
}
