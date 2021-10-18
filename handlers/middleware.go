package handlers

import (
	"backendServer/errors"
	"backendServer/usecases"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CheckAuth() gin.HandlerFunc
}

type MiddlewareImpl struct {
	sessionUseCase usecases.SessionUseCase
}

func CreateMiddleware(sessionUseCase usecases.SessionUseCase) Middleware {
	return &MiddlewareImpl{sessionUseCase: sessionUseCase}
}

func (middleware *MiddlewareImpl) CheckAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := c.Request.Cookie("session_id")
		if err != nil {
			c.String(errors.ResolveErrorToCode(errors.ErrNotAuthorized), errors.ErrNotAuthorized.Error())
			return
		}

		sid := session.Value
		uid, err := middleware.sessionUseCase.GetUID(sid)
		if err != nil {
			session.Expires = time.Now().AddDate(0, 0, -1)
			http.SetCookie(c.Writer, session)
			c.String(errors.ResolveErrorToCode(err), err.Error())
			return
		}

		c.Set("uid", uid)
		c.Set("sid", sid)

		return
	}
}
