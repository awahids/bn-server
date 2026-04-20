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
