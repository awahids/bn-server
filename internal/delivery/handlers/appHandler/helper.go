package apphandler

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	"bn-mobile/server/internal/delivery/data/response"
	"bn-mobile/server/internal/delivery/middleware"

	"github.com/gin-gonic/gin"
)

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func getUserID(c *gin.Context) (string, error) {
	value, exists := c.Get(middleware.ContextUserIDKey)
	if !exists {
		return "", errors.New("missing user context")
	}
	userID, ok := value.(string)
	if !ok || strings.TrimSpace(userID) == "" {
		return "", errors.New("invalid user context")
	}
	return userID, nil
}

func failUnauthorized(c *gin.Context, err error) {
	response.Failed(c, http.StatusUnauthorized, "unauthorized", err.Error())
}

func decodePath(value string) string {
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

func intOrDefault(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func boolOrDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}
