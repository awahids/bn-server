package repointerface

import (
	"context"

	"github.com/awahids/bn-server/internal/domain/models"
)

type AppRepository interface {
	FindUserByID(ctx context.Context, userID string) (*models.User, error)
	FindUserByUsername(ctx context.Context, username string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error

	GetUserProgress(ctx context.Context, userID string, module *string) ([]models.UserProgress, error)
	GetProgressByItem(ctx context.Context, userID, module, itemID string) (*models.UserProgress, error)
	UpsertProgress(ctx context.Context, progress *models.UserProgress) (*models.UserProgress, error)

	GetUserBookmarks(ctx context.Context, userID string, bookmarkType *string) ([]models.Bookmark, error)
	CreateBookmark(ctx context.Context, bookmark *models.Bookmark) (bool, error)
	FindBookmarkByID(ctx context.Context, bookmarkID string) (*models.Bookmark, error)
	DeleteBookmark(ctx context.Context, bookmarkID string) error

	GetDhikrs(ctx context.Context) ([]models.Dhikr, error)
	GetDhikrCountersForDate(ctx context.Context, userID, date string) ([]models.DhikrCounter, error)
	GetDhikrCounter(ctx context.Context, userID, dhikrID, date, session string) (*models.DhikrCounter, error)
	UpsertDhikrCounter(ctx context.Context, counter *models.DhikrCounter) (*models.DhikrCounter, error)

	GetUserQuizAttempts(ctx context.Context, userID string, category *string) ([]models.QuizAttempt, error)
	CreateQuizAttempt(ctx context.Context, attempt *models.QuizAttempt) error
}
