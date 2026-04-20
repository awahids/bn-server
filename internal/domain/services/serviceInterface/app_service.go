package serviceinterface

import (
	"context"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
)

type UpdateUserInput struct {
	Name          *string
	Username      *string
	Streak        *int
	DailyProgress *int
	Preferences   map[string]any
}

type UpsertProgressInput struct {
	Module    string
	ItemID    string
	Progress  int
	Completed bool
	Score     int
	TimeSpent int
}

type CreateBookmarkInput struct {
	Type      string
	ContentID string
	Note      *string
}

type UpsertDhikrCounterInput struct {
	DhikrID   string
	Count     int
	Target    int
	Date      string
	Session   string
	Completed bool
}

type QuizAnswerInput struct {
	QuestionID    string `json:"questionId"`
	UserAnswer    string `json:"userAnswer"`
	CorrectAnswer string `json:"correctAnswer"`
	IsCorrect     bool   `json:"isCorrect"`
	TimeSpent     *int   `json:"timeSpent,omitempty"`
}

type CreateQuizAttemptInput struct {
	Category       string
	Score          int
	TotalQuestions int
	TimeSpent      int
	Answers        []QuizAnswerInput
}

type QuizOverallStats struct {
	TotalAttempts       int     `json:"totalAttempts"`
	AverageScore        float64 `json:"averageScore"`
	BestScore           int     `json:"bestScore"`
	TotalTimeSpent      int     `json:"totalTimeSpent"`
	CategoriesAttempted int     `json:"categoriesAttempted"`
}

type QuizCategoryStats struct {
	Attempts       int        `json:"attempts"`
	AverageScore   float64    `json:"averageScore"`
	BestScore      int        `json:"bestScore"`
	TotalTimeSpent int        `json:"totalTimeSpent"`
	LastAttempt    *time.Time `json:"lastAttempt"`
}

type QuizCategoryBreakdownItem struct {
	Attempts     int     `json:"attempts"`
	TotalScore   int     `json:"totalScore"`
	BestScore    int     `json:"bestScore"`
	TotalTime    int     `json:"totalTime"`
	AverageScore float64 `json:"averageScore"`
	AverageTime  float64 `json:"averageTime"`
}

type QuizRecentPerformance struct {
	AverageScore float64 `json:"averageScore"`
	Trend        string  `json:"trend"`
}

type QuizStatsResponse struct {
	Overall           QuizOverallStats                     `json:"overall"`
	CategoryBreakdown map[string]QuizCategoryBreakdownItem `json:"categoryBreakdown"`
	RecentPerformance *QuizRecentPerformance               `json:"recentPerformance"`
	LastAttempt       *time.Time                           `json:"lastAttempt"`
}

type AppService interface {
	GetUserProfile(ctx context.Context, userID string) (*models.User, error)
	UpdateUserProfile(ctx context.Context, userID string, input UpdateUserInput) (*models.User, error)

	GetProgress(ctx context.Context, userID string, module *string) ([]models.UserProgress, error)
	GetProgressItem(ctx context.Context, userID, module, itemID string) (*models.UserProgress, error)
	UpsertProgress(ctx context.Context, userID string, input UpsertProgressInput) (*models.UserProgress, error)

	GetBookmarks(ctx context.Context, userID string, bookmarkType *string) ([]models.Bookmark, error)
	CreateBookmark(ctx context.Context, userID string, input CreateBookmarkInput) (*models.Bookmark, error)
	DeleteBookmark(ctx context.Context, userID, bookmarkID string) error

	GetDhikrCounters(ctx context.Context, userID, date string) ([]models.DhikrCounter, error)
	UpsertDhikrCounter(ctx context.Context, userID string, input UpsertDhikrCounterInput) (*models.DhikrCounter, error)

	GetQuizAttempts(ctx context.Context, userID string, category *string) ([]models.QuizAttempt, error)
	CreateQuizAttempt(ctx context.Context, userID string, input CreateQuizAttemptInput) (*models.QuizAttempt, error)
	GetQuizCategoryStats(ctx context.Context, userID, category string) (QuizCategoryStats, error)
	GetQuizStats(ctx context.Context, userID string) (QuizStatsResponse, error)
}
