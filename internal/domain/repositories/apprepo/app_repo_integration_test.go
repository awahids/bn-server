package apprepo

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestCreateBookmark_OnConflictIsAtomic(t *testing.T) {
	db := setupIntegrationDB(t)
	tx := beginTx(t, db)
	repo := NewAppRepository(tx)

	user := createTestUser(t, tx)
	note := "first note"
	ctx := context.Background()

	created, err := repo.CreateBookmark(ctx, &models.Bookmark{
		UserID:    user.ID,
		Type:      string(models.BookmarkTypeQuran),
		ContentID: "2:255",
		Note:      &note,
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !created {
		t.Fatal("expected first insert to create bookmark")
	}

	created, err = repo.CreateBookmark(ctx, &models.Bookmark{
		UserID:    user.ID,
		Type:      string(models.BookmarkTypeQuran),
		ContentID: "2:255",
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("expected no error on duplicate insert, got %v", err)
	}
	if created {
		t.Fatal("expected duplicate insert to be ignored")
	}

	var count int64
	if err := tx.WithContext(ctx).
		Model(&models.Bookmark{}).
		Where("user_id = ? AND type = ? AND content_id = ?", user.ID, string(models.BookmarkTypeQuran), "2:255").
		Count(&count).Error; err != nil {
		t.Fatalf("failed counting bookmarks: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 bookmark row, got %d", count)
	}
}

func TestUpsertDhikrCounter_UsesConflictUpdate(t *testing.T) {
	db := setupIntegrationDB(t)
	tx := beginTx(t, db)
	repo := NewAppRepository(tx)

	user := createTestUser(t, tx)
	ctx := context.Background()

	first, err := repo.UpsertDhikrCounter(ctx, &models.DhikrCounter{
		UserID:    user.ID,
		DhikrID:   "morning-dhikr-1",
		Count:     10,
		Target:    33,
		Date:      "2026-04-20",
		Session:   string(models.DhikrSessionMorning),
		Completed: false,
	})
	if err != nil {
		t.Fatalf("expected no error on first upsert, got %v", err)
	}
	if strings.TrimSpace(first.ID) == "" {
		t.Fatal("expected first upsert to return row ID")
	}

	second, err := repo.UpsertDhikrCounter(ctx, &models.DhikrCounter{
		UserID:    user.ID,
		DhikrID:   "morning-dhikr-1",
		Count:     33,
		Target:    33,
		Date:      "2026-04-20",
		Session:   string(models.DhikrSessionMorning),
		Completed: true,
	})
	if err != nil {
		t.Fatalf("expected no error on second upsert, got %v", err)
	}
	if strings.TrimSpace(second.ID) == "" {
		t.Fatal("expected second upsert to return row ID")
	}

	var counters []models.DhikrCounter
	if err := tx.WithContext(ctx).
		Where("user_id = ? AND dhikr_id = ? AND date = ? AND session = ?", user.ID, "morning-dhikr-1", "2026-04-20", string(models.DhikrSessionMorning)).
		Find(&counters).Error; err != nil {
		t.Fatalf("failed querying counters: %v", err)
	}
	if len(counters) != 1 {
		t.Fatalf("expected exactly 1 counter row, got %d", len(counters))
	}
	if counters[0].Count != 33 || !counters[0].Completed {
		t.Fatalf("expected updated counter values (count=33, completed=true), got count=%d completed=%v", counters[0].Count, counters[0].Completed)
	}
}

func setupIntegrationDB(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := strings.TrimSpace(os.Getenv("TEST_POSTGRES_DSN"))
	if dsn == "" {
		t.Skip("TEST_POSTGRES_DSN is not set; skipping postgres integration tests")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed connecting TEST_POSTGRES_DSN: %v", err)
	}

	migrations := []string{
		"000001_create_auth_tables.up.sql",
		"000002_create_app_tables.up.sql",
	}
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed resolving test file path for migrations")
	}
	rootDir := filepath.Clean(filepath.Join(filepath.Dir(currentFile), "../../../../"))
	migrationDir := filepath.Join(rootDir, "internal", "infrastructure", "database", "migrations")

	for _, path := range migrations {
		fullPath := filepath.Join(migrationDir, path)
		sql, err := os.ReadFile(fullPath)
		if err != nil {
			t.Fatalf("failed reading migration %s: %v", fullPath, err)
		}
		if err := db.Exec(string(sql)).Error; err != nil {
			t.Fatalf("failed executing migration %s: %v", fullPath, err)
		}
	}

	return db
}

func beginTx(t *testing.T, db *gorm.DB) *gorm.DB {
	t.Helper()

	tx := db.Begin()
	if tx.Error != nil {
		t.Fatalf("failed to begin transaction: %v", tx.Error)
	}
	t.Cleanup(func() {
		_ = tx.Rollback().Error
	})
	return tx
}

func createTestUser(t *testing.T, db *gorm.DB) *models.User {
	t.Helper()

	suffix := time.Now().UnixNano()
	user := &models.User{
		GoogleID: fmt.Sprintf("google-%d", suffix),
		Email:    fmt.Sprintf("integration-%d@example.com", suffix),
		Name:     "Integration Test User",
		Role:     models.RoleUser,
	}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("failed creating user: %v", err)
	}

	return user
}
