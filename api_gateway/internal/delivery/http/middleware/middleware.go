package middleware

import (
	"api_gateway/pkg/logger"

	"github.com/gin-gonic/gin"
)

// middleware to log request and response
func LogMiddleware(logger logger.Interface) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger.Info("Request: " + c.Request.Method + " " + c.Request.URL.Path)
		c.Next()
		logger.Info("Response: " + c.Request.Method + " " + c.Request.URL.Path)
	}
}
