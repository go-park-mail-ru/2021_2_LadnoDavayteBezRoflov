package handlers

import (
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/logger"
	"expvar"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

type CommonMiddleware interface {
	Logger() gin.HandlerFunc
}

type CommonMiddlewareImpl struct {
	logger logger.Logger
	hits   *expvar.Map
}

func CreateCommonMiddleware(logger logger.Logger, hits *expvar.Map) CommonMiddleware {
	return &CommonMiddlewareImpl{logger: logger, hits: hits}
}

func (middleware *CommonMiddlewareImpl) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		middleware.hits.Add(path, 1)
		if raw != "" {
			path = strings.Join([]string{path, "?", raw}, "; ")
		}

		c.Next()

		if len(c.Errors) > 0 {
			err := customErrors.FindError(c.Errors.Last())
			c.JSON(customErrors.ResolveErrorToCode(err), gin.H{"error": err.Error()})
		}

		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		middleware.logger.Infof("\t%s %s - \"%s %s %d %s\"",
			requestID,
			clientIP,
			method,
			path, // полный путь, на который пришел запрос
			statusCode,
			latency)

		if len(c.Errors) > 0 {
			errorsLog := strings.Join([]string{requestID, ": [", strings.Join(c.Errors.Errors(), "; "), "]"}, "")
			middleware.logger.Error(errorsLog)
		}
	}
}
