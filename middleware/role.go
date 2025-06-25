package middleware

import (
	"food-delivery-app-server/models"
	appErr "food-delivery-app-server/pkg/errors"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func RequireRoles(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			appError := appErr.NewUnauthorized("Missing token claims", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		mapClaims, ok := claims.(jwt.MapClaims)
		if !ok {
			appError := appErr.NewUnauthorized("Invalid token claims", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		role, ok := mapClaims["role"].(string)
		if !ok {
			appError := appErr.NewUnauthorized("Missing user role in token", nil)
			c.JSON(appError.Code, gin.H{"error": appError.Message})
			c.Abort()
			return
		}

		for _, allowed := range roles {
			if strings.EqualFold(role, string(allowed)) {
				c.Next()
				return
			}
		}

		appError := appErr.NewUnauthorized("Access Denied. User Role is not allowed", nil)
		c.JSON(appError.Code, gin.H{"error": appError.Message})
		c.Abort()
	}
}
