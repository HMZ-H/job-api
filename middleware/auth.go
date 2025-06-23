package middleware

import (
	"net/http"
	"strings"
	"job-api/models"
	"job-api/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, models.BaseResponse{
				Success: false,
				Message: "Authorization header required",
				Object:  nil,
				Errors:  []string{"Missing authorization header"},
			})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateJWT(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, models.BaseResponse{
				Success: false,
				Message: "Invalid token",
				Object:  nil,
				Errors:  []string{err.Error()},
			})
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			c.JSON(http.StatusUnauthorized, models.BaseResponse{
				Success: false,
				Message: "User role not found",
				Object:  nil,
			})
			c.Abort()
			return
		}

		if userRole != string(role) {
			c.JSON(http.StatusForbidden, models.BaseResponse{
				Success: false,
				Message: "Insufficient permissions",
				Object:  nil,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
