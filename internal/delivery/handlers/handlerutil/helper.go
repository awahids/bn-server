package handlerutil

import (
	"errors"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/delivery/data/response"
	"github.com/awahids/bn-server/internal/delivery/middleware"
	"github.com/awahids/bn-server/internal/domain/models"

	"github.com/gin-gonic/gin"
)

const RequestTimeout = 8 * time.Second

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
var dateRegex = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

func GetUserID(c *gin.Context) (string, error) {
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

func FailUnauthorized(c *gin.Context, err error) {
	response.Failed(c, http.StatusUnauthorized, "unauthorized", err.Error())
}

func DecodePath(value string) string {
	decoded, err := url.PathUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}

func IntOrDefault(value *int, fallback int) int {
	if value == nil {
		return fallback
	}
	return *value
}

func BoolOrDefault(value *bool, fallback bool) bool {
	if value == nil {
		return fallback
	}
	return *value
}

func IsValidUsername(username string) bool {
	return usernameRegex.MatchString(username)
}

func IsValidDate(value string) bool {
	if !dateRegex.MatchString(value) {
		return false
	}

	parsedDate, err := time.Parse("2006-01-02", value)
	if err != nil {
		return false
	}
	return parsedDate.Format("2006-01-02") == value
}

func IsValidProgressModule(module string) bool {
	switch module {
	case string(models.ModuleHijaiyah), string(models.ModuleQuran), string(models.ModuleDhikr), string(models.ModuleQuiz):
		return true
	default:
		return false
	}
}

func IsValidBookmarkType(bookmarkType string) bool {
	switch bookmarkType {
	case string(models.BookmarkTypeQuran), string(models.BookmarkTypeDhikr):
		return true
	default:
		return false
	}
}

func IsValidDhikrSession(session string) bool {
	switch session {
	case string(models.DhikrSessionMorning), string(models.DhikrSessionEvening):
		return true
	default:
		return false
	}
}
