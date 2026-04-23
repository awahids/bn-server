package pushservice

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/awahids/bn-server/configs"
	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/repositories/repointerface"
)

type PushService struct {
	repo            repointerface.AppRepository
	enabled         bool
	vapidPublicKey  string
	vapidPrivateKey string
	vapidSubject    string
	interval        time.Duration
	httpClient      *http.Client
}

type reminderPayload struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	URL   string `json:"url"`
	Tag   string `json:"tag"`
}

func NewPushService(repo repointerface.AppRepository, cfg configs.PushConfig) *PushService {
	return &PushService{
		repo:            repo,
		enabled:         cfg.Enabled && strings.TrimSpace(cfg.VAPIDPublicKey) != "" && strings.TrimSpace(cfg.VAPIDPrivateKey) != "",
		vapidPublicKey:  strings.TrimSpace(cfg.VAPIDPublicKey),
		vapidPrivateKey: strings.TrimSpace(cfg.VAPIDPrivateKey),
		vapidSubject:    strings.TrimSpace(cfg.VAPIDSubject),
		interval:        cfg.DispatchInterval,
		httpClient:      &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *PushService) IsEnabled() bool {
	return s.enabled
}

func (s *PushService) Run(ctx context.Context) {
	if !s.enabled {
		log.Println("push scheduler disabled")
		return
	}

	if s.interval <= 0 {
		s.interval = time.Minute
	}

	log.Printf("push scheduler started (interval=%s)", s.interval.String())
	// Trigger once shortly after startup so reminders do not wait one full interval.
	s.runDispatchWithTimeout(ctx, 12*time.Second, time.Now())

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("push scheduler stopped")
			return
		case now := <-ticker.C:
			s.runDispatchWithTimeout(ctx, 12*time.Second, now)
		}
	}
}

func (s *PushService) runDispatchWithTimeout(parent context.Context, timeout time.Duration, now time.Time) {
	runCtx, cancel := context.WithTimeout(parent, timeout)
	defer cancel()

	if err := s.DispatchOnce(runCtx, now); err != nil {
		log.Printf("push scheduler dispatch error: %v", err)
	}
}

func (s *PushService) DispatchOnce(ctx context.Context, now time.Time) error {
	timezones, err := s.repo.GetActivePushTimezones(ctx)
	if err != nil {
		return err
	}

	for _, timezone := range timezones {
		timezone = strings.TrimSpace(timezone)
		if timezone == "" {
			continue
		}

		loc, err := time.LoadLocation(timezone)
		if err != nil {
			continue
		}

		localNow := now.In(loc)
		reminderTime := localNow.Format("15:04")
		localDate := localNow.Format("2006-01-02")

		targets, err := s.repo.GetDuePushReminderTargets(ctx, timezone, reminderTime, localDate)
		if err != nil {
			log.Printf("push scheduler query failed timezone=%s: %v", timezone, err)
			continue
		}

		for _, target := range targets {
			if err := s.sendReminder(ctx, target); err != nil {
				log.Printf("push send failed user=%s habit=%s endpoint=%s: %v", target.UserID, target.HabitID, target.Endpoint, err)
			}
		}
	}

	return nil
}

func (s *PushService) sendReminder(ctx context.Context, target models.PushReminderTarget) error {
	payload := reminderPayload{
		Title: "Pengingat Habit",
		Body:  "Saatnya: " + target.HabitName,
		URL:   "/habits",
		Tag:   "habit-reminder-" + target.HabitID,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	subscription := &webpush.Subscription{
		Endpoint: target.Endpoint,
		Keys: webpush.Keys{
			P256dh: target.P256DH,
			Auth:   target.Auth,
		},
	}

	resp, err := webpush.SendNotification(payloadBytes, subscription, &webpush.Options{
		HTTPClient:      s.httpClient,
		Subscriber:      s.vapidSubject,
		VAPIDPublicKey:  s.vapidPublicKey,
		VAPIDPrivateKey: s.vapidPrivateKey,
		TTL:             120,
	})
	if err != nil {
		return err
	}
	if resp == nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusGone || resp.StatusCode == http.StatusNotFound {
		if deleteErr := s.repo.DeletePushSubscriptionByEndpoint(ctx, target.Endpoint); deleteErr != nil {
			log.Printf("failed deleting stale subscription endpoint=%s: %v", target.Endpoint, deleteErr)
		}
	}

	return nil
}
