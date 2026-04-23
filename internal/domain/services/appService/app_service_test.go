package appservice

import (
	"context"
	"errors"
	"testing"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
)

type mockAppRepo struct {
	createBookmarkCreated bool
	createBookmarkErr     error
	lastBookmark          *models.Bookmark
}

func (m *mockAppRepo) FindUserByID(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (m *mockAppRepo) FindUserByUsername(context.Context, string) (*models.User, error) {
	return nil, nil
}

func (m *mockAppRepo) UpdateUser(context.Context, *models.User) error {
	return nil
}

func (m *mockAppRepo) GetUserProgress(context.Context, string, *string) ([]models.UserProgress, error) {
	return nil, nil
}

func (m *mockAppRepo) GetProgressByItem(context.Context, string, string, string) (*models.UserProgress, error) {
	return nil, nil
}

func (m *mockAppRepo) UpsertProgress(context.Context, *models.UserProgress) (*models.UserProgress, error) {
	return nil, nil
}

func (m *mockAppRepo) GetUserHabits(context.Context, string) ([]models.Habit, error) {
	return nil, nil
}

func (m *mockAppRepo) GetUserHabitCompletions(context.Context, string) ([]models.HabitCompletion, error) {
	return nil, nil
}

func (m *mockAppRepo) FindHabitByID(context.Context, string, string) (*models.Habit, error) {
	return nil, nil
}

func (m *mockAppRepo) CreateHabit(context.Context, *models.Habit) error {
	return nil
}

func (m *mockAppRepo) UpdateHabit(context.Context, *models.Habit) error {
	return nil
}

func (m *mockAppRepo) DeleteHabit(context.Context, string, string) error {
	return nil
}

func (m *mockAppRepo) UpsertHabitCompletion(context.Context, *models.HabitCompletion) (*models.HabitCompletion, error) {
	return nil, nil
}

func (m *mockAppRepo) DeleteHabitCompletion(context.Context, string, string, string) error {
	return nil
}

func (m *mockAppRepo) UpsertPushSubscription(context.Context, *models.PushSubscription) error {
	return nil
}

func (m *mockAppRepo) DeletePushSubscriptionByUserEndpoint(context.Context, string, string) error {
	return nil
}

func (m *mockAppRepo) DeletePushSubscriptionByEndpoint(context.Context, string) error {
	return nil
}

func (m *mockAppRepo) GetActivePushTimezones(context.Context) ([]string, error) {
	return nil, nil
}

func (m *mockAppRepo) GetDuePushReminderTargets(context.Context, string, string, string) ([]models.PushReminderTarget, error) {
	return nil, nil
}

func (m *mockAppRepo) GetSchools(context.Context) ([]models.School, error) {
	return nil, nil
}

func (m *mockAppRepo) CreateSchool(context.Context, *models.School) error {
	return nil
}

func (m *mockAppRepo) GetUserBookmarks(context.Context, string, *string) ([]models.Bookmark, error) {
	return nil, nil
}

func (m *mockAppRepo) CreateBookmark(_ context.Context, bookmark *models.Bookmark) (bool, error) {
	m.lastBookmark = bookmark
	return m.createBookmarkCreated, m.createBookmarkErr
}

func (m *mockAppRepo) FindBookmarkByID(context.Context, string) (*models.Bookmark, error) {
	return nil, nil
}

func (m *mockAppRepo) DeleteBookmark(context.Context, string) error {
	return nil
}

func (m *mockAppRepo) GetDhikrs(context.Context) ([]models.Dhikr, error) {
	return nil, nil
}

func (m *mockAppRepo) GetDhikrCountersForDate(context.Context, string, string) ([]models.DhikrCounter, error) {
	return nil, nil
}

func (m *mockAppRepo) GetDhikrCounter(context.Context, string, string, string, string) (*models.DhikrCounter, error) {
	return nil, nil
}

func (m *mockAppRepo) UpsertDhikrCounter(context.Context, *models.DhikrCounter) (*models.DhikrCounter, error) {
	return nil, nil
}

func (m *mockAppRepo) GetUserQuizAttempts(context.Context, string, *string) ([]models.QuizAttempt, error) {
	return nil, nil
}

func (m *mockAppRepo) CreateQuizAttempt(context.Context, *models.QuizAttempt) error {
	return nil
}

func TestCreateBookmark_Success(t *testing.T) {
	repo := &mockAppRepo{createBookmarkCreated: true}
	svc := NewAppService(repo)
	note := "my note"

	bookmark, err := svc.CreateBookmark(context.Background(), "user-1", serviceinterface.CreateBookmarkInput{
		Type:      string(models.BookmarkTypeQuran),
		ContentID: "2:255",
		Note:      &note,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if bookmark == nil {
		t.Fatal("expected bookmark, got nil")
	}
	if bookmark.UserID != "user-1" {
		t.Fatalf("expected user ID user-1, got %s", bookmark.UserID)
	}
	if repo.lastBookmark == nil {
		t.Fatal("expected bookmark to be passed to repository")
	}
}

func TestCreateBookmark_AlreadyExists(t *testing.T) {
	repo := &mockAppRepo{createBookmarkCreated: false}
	svc := NewAppService(repo)

	bookmark, err := svc.CreateBookmark(context.Background(), "user-1", serviceinterface.CreateBookmarkInput{
		Type:      string(models.BookmarkTypeQuran),
		ContentID: "2:255",
	})
	if !errors.Is(err, ErrBookmarkExists) {
		t.Fatalf("expected ErrBookmarkExists, got %v", err)
	}
	if bookmark != nil {
		t.Fatal("expected nil bookmark for duplicate")
	}
}

func TestCreateBookmark_RepoError(t *testing.T) {
	repoErr := errors.New("repository failure")
	repo := &mockAppRepo{
		createBookmarkCreated: false,
		createBookmarkErr:     repoErr,
	}
	svc := NewAppService(repo)

	bookmark, err := svc.CreateBookmark(context.Background(), "user-1", serviceinterface.CreateBookmarkInput{
		Type:      string(models.BookmarkTypeQuran),
		ContentID: "2:255",
	})
	if !errors.Is(err, repoErr) {
		t.Fatalf("expected repo error, got %v", err)
	}
	if bookmark != nil {
		t.Fatal("expected nil bookmark when repository fails")
	}
}
