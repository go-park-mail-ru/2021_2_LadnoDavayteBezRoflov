package handlers

import (
	"backendServer/app/usecases"
	"net/http"
	"time"

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
			session.Expires = time.Now().AddDate(0, 0, -1)
			http.SetCookie(c.Writer, session)
			_ = c.Error(err)
			return
		}

		c.Set("uid", uid)
		c.Set("sid", sid)

		return
	}
}
