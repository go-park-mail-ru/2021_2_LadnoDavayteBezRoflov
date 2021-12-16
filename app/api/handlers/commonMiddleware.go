package handlers

import (
	customErrors "backendServer/pkg/errors"
	"backendServer/pkg/logger"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/penglongli/gin-metrics/ginmetrics"

	"github.com/google/uuid"

	"github.com/gin-gonic/gin"
)

var monitorMutex sync.Mutex

type CommonMiddleware interface {
	Logger() gin.HandlerFunc
}

type CommonMiddlewareImpl struct {
	logger       logger.Logger
	metricsMutex sync.Mutex
	errorMetric  *ginmetrics.Metric
}

func CreateCommonMiddleware(logger logger.Logger) CommonMiddleware {
	monitorMutex.Lock()
	commonMiddlewareImpl := &CommonMiddlewareImpl{logger: logger, errorMetric: ginmetrics.GetMonitor().GetMetric("api_errors")}
	monitorMutex.Unlock()
	return commonMiddlewareImpl
}

func (middleware *CommonMiddlewareImpl) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := uuid.NewString()
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		if raw != "" {
			path = strings.Join([]string{path, "?", raw}, "; ")
		}

		c.Next()

		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				middleware.metricsMutex.Lock()
				_ = middleware.errorMetric.Inc([]string{strconv.Itoa(customErrors.ResolveErrorToCode(err)), err.Error()})
				middleware.metricsMutex.Unlock()
			}

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
