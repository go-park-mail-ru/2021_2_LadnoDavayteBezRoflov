package handlers

import (
	"backendServer/app/usecases"
	"backendServer/pkg/errors"
	"backendServer/pkg/logger"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type Middleware interface {
	CheckAuth() gin.HandlerFunc
	Logger() gin.HandlerFunc
}

type MiddlewareImpl struct {
	sessionUseCase usecases.SessionUseCase
	logger         logger.Logger
}

func CreateMiddleware(sessionUseCase usecases.SessionUseCase, logger logger.Logger) Middleware {
	return &MiddlewareImpl{sessionUseCase: sessionUseCase, logger: logger}
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

func (middleware *MiddlewareImpl) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}

		c.Next()

		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		middleware.logger.Infof("\t%s %s - \"%s %s %d %s\"", requestID, clientIP, method, path, statusCode, latency)

		for _, err := range c.Errors {
			middleware.logger.Errorf("%s - %s", requestID, err.Error())
			c.JSON(errors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		}
	}
}
