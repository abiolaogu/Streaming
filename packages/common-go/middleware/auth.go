package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/streamverse/common-go/errors"
)

// AuthMiddleware validates JWT tokens
func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.NewUnauthorizedError("Unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Invalid token"))
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, errors.NewUnauthorizedError("Invalid token claims"))
			c.Abort()
			return
		}

		// Set user ID in context
		if userID, ok := claims["user_id"].(string); ok {
			c.Set("user_id", userID)
		}

		if email, ok := claims["email"].(string); ok {
			c.Set("email", email)
		}

		c.Next()
	}
}

// RequireRole checks if user has required role
func RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roles, exists := c.Get("roles")
		if !exists {
			c.JSON(http.StatusForbidden, errors.NewForbiddenError("No roles found"))
			c.Abort()
			return
		}

		rolesList, ok := roles.([]string)
		if !ok {
			c.JSON(http.StatusForbidden, errors.NewForbiddenError("Invalid roles format"))
			c.Abort()
			return
		}

		for _, role := range rolesList {
			if role == requiredRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, errors.NewForbiddenError("Insufficient permissions"))
		c.Abort()
	}
}

