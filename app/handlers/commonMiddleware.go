package handlers

import (
	"backendServer/pkg/errors"
	"backendServer/pkg/logger"
	"fmt"
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
}

func CreateCommonMiddleware(logger logger.Logger) CommonMiddleware {
	return &CommonMiddlewareImpl{logger: logger}
}

func (middleware *CommonMiddlewareImpl) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		var resultPath strings.Builder
		resultPath.WriteString(path)
		if raw != "" {
			resultPath.WriteString("?")
			resultPath.WriteString(raw)
		}

		c.Next()

		timeStamp := time.Now()
		latency := timeStamp.Sub(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		middleware.logger.Infof("\t%s %s - \"%s %s %d %s\"", requestID, clientIP, method, resultPath.String(), statusCode, latency)

		if len(c.Errors) > 0 {
			var errorsLogs strings.Builder
			errorsLogs.WriteString(requestID)
			for _, err := range c.Errors {
				_, printErr := fmt.Fprintf(&errorsLogs, "\n\t\t\t\t\t\t\t\t\t%s", err.Error())
				if printErr != nil {
					middleware.logger.Error(printErr)
					return
				}

			}
			middleware.logger.Error(errorsLogs.String())
			c.JSON(errors.ResolveErrorToCode(c.Errors.Last()), gin.H{"error": c.Errors.Last().Error()})
		}
	}
}
