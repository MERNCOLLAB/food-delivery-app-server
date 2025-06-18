package middleware

import (
	appErr "food-delivery-app-server/pkg/errors" // alias to avoid name clash
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			log.Printf("Error: %v", err)

			if appError, ok := err.(*appErr.AppError); ok {
				// Custom Errors
				c.JSON(appError.Code, gin.H{
					"success": false,
					"message": appError.Message,
					"error":   appError.Err.Error(),
				})
				return
			}

			// Fallback Error
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "Internal server error",
				"error":   err.Error(),
			})
		}
	}
}
