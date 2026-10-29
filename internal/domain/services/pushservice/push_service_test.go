package pushservice

import (
	"context"
	"testing"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/repositories/repointerface"
)

type fakeRepo struct {
	repointerface.AppRepository
	timezones      []string
	dailyCalls     []time.Time
	dailyTimezones []string
}

func (f *fakeRepo) GetActivePushTimezones(context.Context) ([]string, error) {
	return f.timezones, nil
}

func (f *fakeRepo) GetDuePushReminderTargets(context.Context, string, string, string) ([]models.PushReminderTarget, error) {
	return nil, nil
}

func (f *fakeRepo) GetDueDailyReminderTargets(_ context.Context, timezone string, dayStart time.Time) ([]models.PushReminderTarget, error) {
	f.dailyTimezones = append(f.dailyTimezones, timezone)
	f.dailyCalls = append(f.dailyCalls, dayStart)
	return nil, nil
}

func TestDispatchOnceQueriesDailyReminderAtLocalTime(t *testing.T) {
	repo := &fakeRepo{timezones: []string{"Asia/Jakarta", "Asia/Jayapura"}}
	service := &PushService{repo: repo}

	// 12:00 UTC is 19:00 in Jakarta (UTC+7) and 21:00 in Jayapura (UTC+9).
	now := time.Date(2026, 10, 29, 12, 0, 0, 0, time.UTC)
	if err := service.DispatchOnce(context.Background(), now); err != nil {
		t.Fatalf("DispatchOnce: %v", err)
	}

	if len(repo.dailyTimezones) != 1 || repo.dailyTimezones[0] != "Asia/Jakarta" {
		t.Fatalf("daily reminder timezones = %v, want [Asia/Jakarta]", repo.dailyTimezones)
	}
	jakarta, _ := time.LoadLocation("Asia/Jakarta")
	wantStart := time.Date(2026, 10, 29, 0, 0, 0, 0, jakarta)
	if !repo.dailyCalls[0].Equal(wantStart) {
		t.Fatalf("dayStart = %s, want %s", repo.dailyCalls[0], wantStart)
	}
}

func TestDispatchOnceSkipsDailyReminderOutsideItsMinute(t *testing.T) {
	repo := &fakeRepo{timezones: []string{"Asia/Jakarta"}}
	service := &PushService{repo: repo}

	for _, minute := range []int{59, 1} {
		hour := 11
		if minute == 1 {
			hour = 12
		}
		now := time.Date(2026, 10, 29, hour, minute, 0, 0, time.UTC)
		if err := service.DispatchOnce(context.Background(), now); err != nil {
			t.Fatalf("DispatchOnce: %v", err)
		}
	}

	if len(repo.dailyCalls) != 0 {
		t.Fatalf("daily reminder queried %d times outside 19:00, want 0", len(repo.dailyCalls))
	}
}

func TestDailyReminderPayloadIsTaggedPerDay(t *testing.T) {
	payload := dailyReminderPayload("2026-10-29")
	if payload.Tag != "daily-reminder-2026-10-29" || payload.URL != "/modules" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
}
