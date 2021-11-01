package handlers

import (
	customErrors "backendServer/pkg/errors"

	"github.com/gin-gonic/gin"
)

func NoRouteHandler(c *gin.Context) {
	c.JSON(customErrors.ResolveErrorToCode(customErrors.ErrNotImplemented), gin.H{"error": customErrors.ErrNotImplemented.Error()})
}
