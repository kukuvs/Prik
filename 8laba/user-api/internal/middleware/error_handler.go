package middleware

import (
	"log"
	"net/http"
	"user-api/internal/models"

	"github.com/gin-gonic/gin"
)

// ErrorHandler middleware для централизованной обработки ошибок
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			log.Printf("Error: %v", err)

			c.JSON(http.StatusInternalServerError, models.ErrorResponse{
				Error:   "Internal server error",
				Message: err.Error(),
			})
		}
	}
}

// Logger middleware для логирования запросов
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("[%s] %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// CORS middleware для разрешения кросс-доменных запросов
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
