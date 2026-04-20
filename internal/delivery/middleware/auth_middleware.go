package middleware

import (
	"net/http"
	"strings"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserIDKey = "userID"
	ContextUserRole  = "userRole"
	ContextUserEmail = "userEmail"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Failed(c, http.StatusUnauthorized, "missing authorization header", "missing token")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			response.Failed(c, http.StatusUnauthorized, "invalid authorization header", "invalid token format")
			c.Abort()
			return
		}

		claims, err := utils.ParseAccessToken(jwtSecret, strings.TrimSpace(parts[1]))
		if err != nil {
			response.Failed(c, http.StatusUnauthorized, "invalid access token", err.Error())
			c.Abort()
			return
		}

		c.Set(ContextUserIDKey, claims.Subject)
		c.Set(ContextUserRole, claims.Role)
		c.Set(ContextUserEmail, claims.Email)
		c.Next()
	}
}
