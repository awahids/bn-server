package apprepo

import (
	"context"
	"errors"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/repositories/repointerface"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type appRepository struct {
	db *gorm.DB
}

func NewAppRepository(db *gorm.DB) repointerface.AppRepository {
	return &appRepository{db: db}
}

func (r *appRepository) FindUserByID(ctx context.Context, userID string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).
		Preload("Progress").
		Preload("Bookmarks").
		Preload("Habits").
		Preload("HabitLogs").
		Preload("DhikrCounters").
		Preload("QuizAttempts").
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *appRepository) FindUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *appRepository) UpdateUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

func (r *appRepository) GetUserProgress(ctx context.Context, userID string, module *string) ([]models.UserProgress, error) {
	var progress []models.UserProgress
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if module != nil && *module != "" {
		query = query.Where("module = ?", *module)
	}

	err := query.Order("last_accessed DESC").Find(&progress).Error
	if err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *appRepository) GetProgressByItem(ctx context.Context, userID, module, itemID string) (*models.UserProgress, error) {
	var progress models.UserProgress
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND module = ? AND item_id = ?", userID, module, itemID).
		First(&progress).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &progress, nil
}

func (r *appRepository) UpsertProgress(ctx context.Context, progress *models.UserProgress) (*models.UserProgress, error) {
	progress.LastAccessed = time.Now()
	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_id"},
					{Name: "module"},
					{Name: "item_id"},
				},
				DoUpdates: clause.Assignments(map[string]any{
					"progress":      progress.Progress,
					"completed":     progress.Completed,
					"score":         progress.Score,
					"time_spent":    progress.TimeSpent,
					"last_accessed": progress.LastAccessed,
					"updated_at":    time.Now(),
					// Ensure soft-deleted rows can be revived by upsert.
					"deleted_at": nil,
				}),
			},
			clause.Returning{},
		).
		Create(progress).Error; err != nil {
		return nil, err
	}
	return progress, nil
}

func (r *appRepository) GetUserHabits(ctx context.Context, userID string) ([]models.Habit, error) {
	var habits []models.Habit
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_at ASC").
		Find(&habits).Error
	if err != nil {
		return nil, err
	}
	return habits, nil
}

func (r *appRepository) GetUserHabitCompletions(ctx context.Context, userID string) ([]models.HabitCompletion, error) {
	var completions []models.HabitCompletion
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND completed = ?", userID, true).
		Order("date ASC, created_at ASC").
		Find(&completions).Error
	if err != nil {
		return nil, err
	}
	return completions, nil
}

func (r *appRepository) FindHabitByID(ctx context.Context, userID, habitID string) (*models.Habit, error) {
	var habit models.Habit
	err := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", habitID, userID).
		First(&habit).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &habit, nil
}

func (r *appRepository) CreateHabit(ctx context.Context, habit *models.Habit) error {
	return r.db.WithContext(ctx).Create(habit).Error
}

func (r *appRepository) UpdateHabit(ctx context.Context, habit *models.Habit) error {
	return r.db.WithContext(ctx).Save(habit).Error
}

func (r *appRepository) DeleteHabit(ctx context.Context, userID, habitID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND user_id = ?", habitID, userID).
		Delete(&models.Habit{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *appRepository) UpsertHabitCompletion(ctx context.Context, completion *models.HabitCompletion) (*models.HabitCompletion, error) {
	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_id"},
					{Name: "habit_id"},
					{Name: "date"},
				},
				DoUpdates: clause.Assignments(map[string]any{
					"completed":  completion.Completed,
					"updated_at": time.Now(),
					"deleted_at": nil,
				}),
			},
			clause.Returning{},
		).
		Create(completion).Error; err != nil {
		return nil, err
	}
	return completion, nil
}

func (r *appRepository) DeleteHabitCompletion(ctx context.Context, userID, habitID, date string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND habit_id = ? AND date = ?", userID, habitID, date).
		Delete(&models.HabitCompletion{}).Error
}

func (r *appRepository) GetUserBookmarks(ctx context.Context, userID string, bookmarkType *string) ([]models.Bookmark, error) {
	var bookmarks []models.Bookmark
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if bookmarkType != nil && *bookmarkType != "" {
		query = query.Where("type = ?", *bookmarkType)
	}

	err := query.Order("created_at DESC").Find(&bookmarks).Error
	if err != nil {
		return nil, err
	}
	return bookmarks, nil
}

func (r *appRepository) CreateBookmark(ctx context.Context, bookmark *models.Bookmark) (bool, error) {
	result := r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "user_id"},
				{Name: "type"},
				{Name: "content_id"},
			},
			DoNothing: true,
		}).
		Create(bookmark)
	if result.Error != nil {
		return false, result.Error
	}

	return result.RowsAffected > 0, nil
}

func (r *appRepository) FindBookmarkByID(ctx context.Context, bookmarkID string) (*models.Bookmark, error) {
	var bookmark models.Bookmark
	err := r.db.WithContext(ctx).Where("id = ?", bookmarkID).First(&bookmark).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bookmark, nil
}

func (r *appRepository) DeleteBookmark(ctx context.Context, bookmarkID string) error {
	result := r.db.WithContext(ctx).Where("id = ?", bookmarkID).Delete(&models.Bookmark{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *appRepository) GetDhikrCountersForDate(ctx context.Context, userID, date string) ([]models.DhikrCounter, error) {
	var counters []models.DhikrCounter
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND date = ?", userID, date).
		Order("dhikr_id ASC").
		Find(&counters).Error
	if err != nil {
		return nil, err
	}
	return counters, nil
}

func (r *appRepository) GetDhikrs(ctx context.Context) ([]models.Dhikr, error) {
	var dhikrs []models.Dhikr
	err := r.db.WithContext(ctx).Order("category ASC, id ASC").Find(&dhikrs).Error
	if err != nil {
		return nil, err
	}
	return dhikrs, nil
}

func (r *appRepository) GetDhikrCounter(ctx context.Context, userID, dhikrID, date, session string) (*models.DhikrCounter, error) {
	var counter models.DhikrCounter
	err := r.db.WithContext(ctx).
		Where("user_id = ? AND dhikr_id = ? AND date = ? AND session = ?", userID, dhikrID, date, session).
		First(&counter).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &counter, nil
}

func (r *appRepository) UpsertDhikrCounter(ctx context.Context, counter *models.DhikrCounter) (*models.DhikrCounter, error) {
	if err := r.db.WithContext(ctx).
		Clauses(
			clause.OnConflict{
				Columns: []clause.Column{
					{Name: "user_id"},
					{Name: "dhikr_id"},
					{Name: "date"},
					{Name: "session"},
				},
				DoUpdates: clause.Assignments(map[string]any{
					"count":      counter.Count,
					"target":     counter.Target,
					"completed":  counter.Completed,
					"updated_at": time.Now(),
					// Ensure soft-deleted rows can be revived by upsert.
					"deleted_at": nil,
				}),
			},
			clause.Returning{},
		).
		Create(counter).Error; err != nil {
		return nil, err
	}

	return counter, nil
}

func (r *appRepository) GetUserQuizAttempts(ctx context.Context, userID string, category *string) ([]models.QuizAttempt, error) {
	var attempts []models.QuizAttempt
	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if category != nil && *category != "" {
		query = query.Where("category = ?", *category)
	}

	err := query.Order("completed_at DESC").Find(&attempts).Error
	if err != nil {
		return nil, err
	}
	return attempts, nil
}

func (r *appRepository) CreateQuizAttempt(ctx context.Context, attempt *models.QuizAttempt) error {
	return r.db.WithContext(ctx).Create(attempt).Error
}
