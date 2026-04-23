package appservice

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"sort"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/domain/models"
	"github.com/awahids/bn-server/internal/domain/repositories/repointerface"
	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrHabitNotFound     = errors.New("habit not found")
	ErrHabitInvalidData  = errors.New("invalid habit data")
	ErrPushInvalidData   = errors.New("invalid push subscription data")
	ErrSchoolInvalidData = errors.New("invalid school data")
	ErrBookmarkExists    = errors.New("bookmark already exists for this content")
	ErrBookmarkNotFound  = errors.New("bookmark not found")
	ErrBookmarkForbidden = errors.New("forbidden: you can only delete your own bookmarks")
	ErrScoreMismatch     = errors.New("score does not match the provided answers")
)

var schoolJenjangAliases = map[string]string{
	"tk":      "TK",
	"sd":      "SD",
	"smp":     "SMP",
	"sma":     "SMA",
	"smk":     "SMK",
	"mi":      "MI",
	"mts":     "MTs",
	"ma":      "MA",
	"lainnya": "Lainnya",
}

type appService struct {
	repo repointerface.AppRepository
}

func NewAppService(repo repointerface.AppRepository) serviceinterface.AppService {
	return &appService{repo: repo}
}

func (s *appService) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *appService) UpdateUserProfile(ctx context.Context, userID string, input serviceinterface.UpdateUserInput) (*models.User, error) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if input.Username != nil {
		existing, err := s.repo.FindUserByUsername(ctx, strings.TrimSpace(*input.Username))
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != userID {
			return nil, ErrUsernameTaken
		}
		trimmed := strings.TrimSpace(*input.Username)
		user.Username = &trimmed
	}

	if input.Name != nil {
		user.Name = strings.TrimSpace(*input.Name)
	}
	if input.Streak != nil {
		user.Streak = *input.Streak
	}
	if input.DailyProgress != nil {
		user.DailyProgress = *input.DailyProgress
	}
	if input.Preferences != nil {
		preferencesJSON, err := json.Marshal(input.Preferences)
		if err != nil {
			return nil, err
		}
		user.Preferences = json.RawMessage(preferencesJSON)
	}

	user.LastActive = time.Now()
	if err := s.repo.UpdateUser(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *appService) GetProgress(ctx context.Context, userID string, module *string) ([]models.UserProgress, error) {
	return s.repo.GetUserProgress(ctx, userID, module)
}

func (s *appService) GetProgressItem(ctx context.Context, userID, module, itemID string) (*models.UserProgress, error) {
	return s.repo.GetProgressByItem(ctx, userID, module, itemID)
}

func (s *appService) UpsertProgress(ctx context.Context, userID string, input serviceinterface.UpsertProgressInput) (*models.UserProgress, error) {
	progress := &models.UserProgress{
		UserID:       userID,
		Module:       input.Module,
		ItemID:       input.ItemID,
		Progress:     input.Progress,
		Completed:    input.Completed,
		Score:        input.Score,
		TimeSpent:    input.TimeSpent,
		LastAccessed: time.Now(),
	}
	return s.repo.UpsertProgress(ctx, progress)
}

func (s *appService) GetHabits(ctx context.Context, userID string) ([]models.Habit, error) {
	return s.repo.GetUserHabits(ctx, userID)
}

func (s *appService) GetHabitCompletions(ctx context.Context, userID string) ([]models.HabitCompletion, error) {
	return s.repo.GetUserHabitCompletions(ctx, userID)
}

func (s *appService) CreateHabit(ctx context.Context, userID string, input serviceinterface.CreateHabitInput) (*models.Habit, error) {
	if input.ReminderEnabled && strings.TrimSpace(input.ReminderTime) == "" {
		return nil, ErrHabitInvalidData
	}

	habit := &models.Habit{
		UserID:          userID,
		Name:            input.Name,
		Category:        input.Category,
		ReminderTime:    input.ReminderTime,
		ReminderEnabled: input.ReminderEnabled,
	}
	if err := s.repo.CreateHabit(ctx, habit); err != nil {
		return nil, err
	}
	return habit, nil
}

func (s *appService) UpdateHabit(ctx context.Context, userID, habitID string, input serviceinterface.UpdateHabitInput) (*models.Habit, error) {
	habit, err := s.repo.FindHabitByID(ctx, userID, habitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}

	if input.Name != nil {
		habit.Name = strings.TrimSpace(*input.Name)
	}
	if input.Category != nil {
		habit.Category = strings.TrimSpace(*input.Category)
	}
	if input.ReminderTime != nil {
		habit.ReminderTime = strings.TrimSpace(*input.ReminderTime)
	}
	if input.ReminderEnabled != nil {
		habit.ReminderEnabled = *input.ReminderEnabled
	}

	if !habit.ReminderEnabled {
		habit.ReminderTime = ""
	}
	if habit.ReminderEnabled && strings.TrimSpace(habit.ReminderTime) == "" {
		return nil, ErrHabitInvalidData
	}

	if err := s.repo.UpdateHabit(ctx, habit); err != nil {
		return nil, err
	}
	return habit, nil
}

func (s *appService) DeleteHabit(ctx context.Context, userID, habitID string) error {
	habit, err := s.repo.FindHabitByID(ctx, userID, habitID)
	if err != nil {
		return err
	}
	if habit == nil {
		return ErrHabitNotFound
	}

	if err := s.repo.DeleteHabit(ctx, userID, habitID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrHabitNotFound
		}
		return err
	}
	return nil
}

func (s *appService) SetHabitCompletion(ctx context.Context, userID string, input serviceinterface.SetHabitCompletionInput) (*models.HabitCompletion, error) {
	habit, err := s.repo.FindHabitByID(ctx, userID, input.HabitID)
	if err != nil {
		return nil, err
	}
	if habit == nil {
		return nil, ErrHabitNotFound
	}

	if !input.Completed {
		if err := s.repo.DeleteHabitCompletion(ctx, userID, input.HabitID, input.Date); err != nil {
			return nil, err
		}
		return nil, nil
	}

	completion := &models.HabitCompletion{
		UserID:    userID,
		HabitID:   input.HabitID,
		Date:      input.Date,
		Completed: true,
	}

	return s.repo.UpsertHabitCompletion(ctx, completion)
}

func (s *appService) UpsertPushSubscription(ctx context.Context, userID string, input serviceinterface.UpsertPushSubscriptionInput) error {
	endpoint := strings.TrimSpace(input.Endpoint)
	p256dh := strings.TrimSpace(input.Keys.P256DH)
	auth := strings.TrimSpace(input.Keys.Auth)
	timezone := normalizeTimeZone(input.Timezone)

	if endpoint == "" || p256dh == "" || auth == "" {
		return ErrPushInvalidData
	}
	if len(endpoint) > 2048 || len(p256dh) > 512 || len(auth) > 255 {
		return ErrPushInvalidData
	}

	subscription := &models.PushSubscription{
		UserID:         userID,
		Endpoint:       endpoint,
		P256DH:         p256dh,
		Auth:           auth,
		ExpirationTime: input.ExpirationTime,
		Timezone:       timezone,
		Enabled:        true,
		LastSeenAt:     time.Now(),
	}

	return s.repo.UpsertPushSubscription(ctx, subscription)
}

func (s *appService) DeletePushSubscription(ctx context.Context, userID, endpoint string) error {
	trimmedEndpoint := strings.TrimSpace(endpoint)
	if trimmedEndpoint == "" {
		return ErrPushInvalidData
	}
	if err := s.repo.DeletePushSubscriptionByUserEndpoint(ctx, userID, trimmedEndpoint); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		return err
	}
	return nil
}

func (s *appService) GetSchools(ctx context.Context) ([]models.School, error) {
	return s.repo.GetSchools(ctx)
}

func (s *appService) CreateSchool(ctx context.Context, userID string, input serviceinterface.CreateSchoolInput) (*models.School, error) {
	name := strings.TrimSpace(input.Name)
	location := strings.TrimSpace(input.Location)
	jenjang, jenjangOK := normalizeSchoolJenjang(input.Jenjang)
	statusSekolah, statusSekolahOK := normalizeSchoolStatusSekolah(input.StatusSekolah)
	mapURL := strings.TrimSpace(input.MapURL)
	contact := ""
	description := ""

	if input.Contact != nil {
		contact = strings.TrimSpace(*input.Contact)
	}
	if input.Description != nil {
		description = strings.TrimSpace(*input.Description)
	}

	if name == "" || location == "" || mapURL == "" || input.MonthlyFee < 0 || !jenjangOK || !statusSekolahOK {
		return nil, ErrSchoolInvalidData
	}

	if len(contact) > 100 || len(description) > 1000 {
		return nil, ErrSchoolInvalidData
	}

	school := &models.School{
		UserID:        userID,
		Name:          name,
		Location:      location,
		Jenjang:       jenjang,
		StatusSekolah: statusSekolah,
		MonthlyFee:    input.MonthlyFee,
		MapURL:        mapURL,
		Contact:       contact,
		Description:   description,
	}

	if err := s.repo.CreateSchool(ctx, school); err != nil {
		return nil, err
	}
	return school, nil
}

func normalizeSchoolJenjang(value string) (string, bool) {
	key := strings.ToLower(strings.TrimSpace(value))
	normalized, ok := schoolJenjangAliases[key]
	return normalized, ok
}

func normalizeSchoolStatusSekolah(value string) (string, bool) {
	status := strings.ToLower(strings.TrimSpace(value))
	switch status {
	case "negeri", "swasta":
		return status, true
	default:
		return "", false
	}
}

func normalizeTimeZone(value string) string {
	timezone := strings.TrimSpace(value)
	if timezone == "" {
		return "UTC"
	}
	if _, err := time.LoadLocation(timezone); err != nil {
		return "UTC"
	}
	return timezone
}

func (s *appService) GetBookmarks(ctx context.Context, userID string, bookmarkType *string) ([]models.Bookmark, error) {
	return s.repo.GetUserBookmarks(ctx, userID, bookmarkType)
}

func (s *appService) CreateBookmark(ctx context.Context, userID string, input serviceinterface.CreateBookmarkInput) (*models.Bookmark, error) {
	bookmark := &models.Bookmark{
		UserID:    userID,
		Type:      input.Type,
		ContentID: input.ContentID,
		Note:      input.Note,
		CreatedAt: time.Now(),
	}
	created, err := s.repo.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}
	if !created {
		return nil, ErrBookmarkExists
	}
	return bookmark, nil
}

func (s *appService) DeleteBookmark(ctx context.Context, userID, bookmarkID string) error {
	bookmark, err := s.repo.FindBookmarkByID(ctx, bookmarkID)
	if err != nil {
		return err
	}
	if bookmark == nil {
		return ErrBookmarkNotFound
	}
	if bookmark.UserID != userID {
		return ErrBookmarkForbidden
	}

	if err := s.repo.DeleteBookmark(ctx, bookmarkID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrBookmarkNotFound
		}
		return err
	}
	return nil
}

func (s *appService) GetDhikrs(ctx context.Context) ([]models.Dhikr, error) {
	return s.repo.GetDhikrs(ctx)
}

func (s *appService) GetDhikrCounters(ctx context.Context, userID, date string) ([]models.DhikrCounter, error) {
	return s.repo.GetDhikrCountersForDate(ctx, userID, date)
}

func (s *appService) UpsertDhikrCounter(ctx context.Context, userID string, input serviceinterface.UpsertDhikrCounterInput) (*models.DhikrCounter, error) {
	counter := &models.DhikrCounter{
		UserID:    userID,
		DhikrID:   input.DhikrID,
		Count:     input.Count,
		Target:    input.Target,
		Date:      input.Date,
		Session:   input.Session,
		Completed: input.Completed,
	}
	return s.repo.UpsertDhikrCounter(ctx, counter)
}

func (s *appService) GetQuizAttempts(ctx context.Context, userID string, category *string) ([]models.QuizAttempt, error) {
	return s.repo.GetUserQuizAttempts(ctx, userID, category)
}

func (s *appService) CreateQuizAttempt(ctx context.Context, userID string, input serviceinterface.CreateQuizAttemptInput) (*models.QuizAttempt, error) {
	correctAnswers := 0
	for _, answer := range input.Answers {
		if answer.IsCorrect {
			correctAnswers++
		}
	}
	expectedScore := int(math.Round((float64(correctAnswers) / float64(input.TotalQuestions)) * 100.0))
	if int(math.Abs(float64(input.Score-expectedScore))) > 1 {
		return nil, ErrScoreMismatch
	}

	answersJSON, err := json.Marshal(input.Answers)
	if err != nil {
		return nil, err
	}

	attempt := &models.QuizAttempt{
		UserID:         userID,
		Category:       input.Category,
		Score:          input.Score,
		TotalQuestions: input.TotalQuestions,
		TimeSpent:      input.TimeSpent,
		Answers:        json.RawMessage(answersJSON),
		CompletedAt:    time.Now(),
	}
	if err := s.repo.CreateQuizAttempt(ctx, attempt); err != nil {
		return nil, err
	}
	return attempt, nil
}

func (s *appService) GetQuizCategoryStats(ctx context.Context, userID, category string) (serviceinterface.QuizCategoryStats, error) {
	attempts, err := s.repo.GetUserQuizAttempts(ctx, userID, &category)
	if err != nil {
		return serviceinterface.QuizCategoryStats{}, err
	}

	if len(attempts) == 0 {
		return serviceinterface.QuizCategoryStats{}, nil
	}

	totalScore := 0
	totalTime := 0
	bestScore := attempts[0].Score
	lastAttempt := attempts[0].CompletedAt

	for _, attempt := range attempts {
		totalScore += attempt.Score
		totalTime += attempt.TimeSpent
		if attempt.Score > bestScore {
			bestScore = attempt.Score
		}
		if attempt.CompletedAt.After(lastAttempt) {
			lastAttempt = attempt.CompletedAt
		}
	}

	return serviceinterface.QuizCategoryStats{
		Attempts:       len(attempts),
		AverageScore:   float64(totalScore) / float64(len(attempts)),
		BestScore:      bestScore,
		TotalTimeSpent: totalTime,
		LastAttempt:    &lastAttempt,
	}, nil
}

func (s *appService) GetQuizStats(ctx context.Context, userID string) (serviceinterface.QuizStatsResponse, error) {
	allAttempts, err := s.repo.GetUserQuizAttempts(ctx, userID, nil)
	if err != nil {
		return serviceinterface.QuizStatsResponse{}, err
	}

	if len(allAttempts) == 0 {
		return serviceinterface.QuizStatsResponse{
			Overall:           serviceinterface.QuizOverallStats{},
			CategoryBreakdown: map[string]serviceinterface.QuizCategoryBreakdownItem{},
			RecentPerformance: nil,
			LastAttempt:       nil,
		}, nil
	}

	totalScore := 0
	totalTimeSpent := 0
	bestScore := allAttempts[0].Score
	categorySet := map[string]struct{}{}
	categoryBreakdown := map[string]serviceinterface.QuizCategoryBreakdownItem{}
	lastAttempt := allAttempts[0].CompletedAt

	for _, attempt := range allAttempts {
		totalScore += attempt.Score
		totalTimeSpent += attempt.TimeSpent
		if attempt.Score > bestScore {
			bestScore = attempt.Score
		}
		categorySet[attempt.Category] = struct{}{}
		if attempt.CompletedAt.After(lastAttempt) {
			lastAttempt = attempt.CompletedAt
		}

		item := categoryBreakdown[attempt.Category]
		item.Attempts++
		item.TotalScore += attempt.Score
		if attempt.Score > item.BestScore {
			item.BestScore = attempt.Score
		}
		item.TotalTime += attempt.TimeSpent
		categoryBreakdown[attempt.Category] = item
	}

	for key, item := range categoryBreakdown {
		item.AverageScore = float64(item.TotalScore) / float64(item.Attempts)
		item.AverageTime = float64(item.TotalTime) / float64(item.Attempts)
		categoryBreakdown[key] = item
	}

	sortedAttempts := make([]models.QuizAttempt, len(allAttempts))
	copy(sortedAttempts, allAttempts)
	sort.Slice(sortedAttempts, func(i, j int) bool {
		return sortedAttempts[i].CompletedAt.After(sortedAttempts[j].CompletedAt)
	})
	if len(sortedAttempts) > 10 {
		sortedAttempts = sortedAttempts[:10]
	}

	recentPerformance := buildRecentPerformance(sortedAttempts)

	return serviceinterface.QuizStatsResponse{
		Overall: serviceinterface.QuizOverallStats{
			TotalAttempts:       len(allAttempts),
			AverageScore:        float64(totalScore) / float64(len(allAttempts)),
			BestScore:           bestScore,
			TotalTimeSpent:      totalTimeSpent,
			CategoriesAttempted: len(categorySet),
		},
		CategoryBreakdown: categoryBreakdown,
		RecentPerformance: recentPerformance,
		LastAttempt:       &lastAttempt,
	}, nil
}

func buildRecentPerformance(attempts []models.QuizAttempt) *serviceinterface.QuizRecentPerformance {
	if len(attempts) == 0 {
		return nil
	}

	scores := make([]int, 0, len(attempts))
	total := 0
	for _, attempt := range attempts {
		scores = append(scores, attempt.Score)
		total += attempt.Score
	}

	return &serviceinterface.QuizRecentPerformance{
		AverageScore: float64(total) / float64(len(scores)),
		Trend:        calculateTrend(scores),
	}
}

func calculateTrend(scores []int) string {
	if len(scores) < 3 {
		return "stable"
	}

	half := len(scores) / 2
	firstHalf := scores[:half]
	secondHalf := scores[half:]

	firstAvg := avgInt(firstHalf)
	secondAvg := avgInt(secondHalf)
	diff := secondAvg - firstAvg

	if diff > 5 {
		return "improving"
	}
	if diff < -5 {
		return "declining"
	}
	return "stable"
}

func avgInt(items []int) float64 {
	if len(items) == 0 {
		return 0
	}
	total := 0
	for _, item := range items {
		total += item
	}
	return float64(total) / float64(len(items))
}

func (s *appService) GetAchievements(ctx context.Context, userID string) ([]serviceinterface.AchievementItem, error) {
	// Re-calculating same stats logic as client to award achievements.
	user, err := s.repo.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return nil, ErrUserNotFound
	}
	progress, _ := s.repo.GetUserProgress(ctx, userID, nil)
	quizStats, _ := s.GetQuizStats(ctx, userID)

	hijaiyahCompleted := 0
	for _, p := range progress {
		if p.Module == "hijaiyah" && p.Completed {
			hijaiyahCompleted++
		}
	}

	achievements := []serviceinterface.AchievementItem{
		{
			ID:          "first-letter",
			Title:       "Huruf Pertama",
			Description: "Selesaikan huruf Hijaiyah pertama",
			Icon:        "Languages",
			Unlocked:    hijaiyahCompleted > 0,
		},
		{
			ID:          "week-streak",
			Title:       "Seminggu Berturut",
			Description: "Belajar 7 hari berturut-turut",
			Icon:        "Flame",
			Unlocked:    user.Streak >= 7,
		},
		{
			ID:          "quiz-master",
			Title:       "Master Kuis",
			Description: "Dapatkan skor 90% atau lebih",
			Icon:        "Trophy",
			Unlocked:    quizStats.Overall.BestScore >= 90,
		},
	}

	todayStr := time.Now().Format("2006-01-02")
	dhikrTodayCounters, _ := s.repo.GetDhikrCountersForDate(ctx, userID, todayStr)
	hasMorning := false
	hasEvening := false
	for _, dc := range dhikrTodayCounters {
		if dc.Completed {
			if dc.Session == "morning" {
				hasMorning = true
			}
			if dc.Session == "evening" {
				hasEvening = true
			}
		}
	}

	achievements = append(achievements, serviceinterface.AchievementItem{
		ID:          "dhikr-complete",
		Title:       "Dhikr Lengkap",
		Description: "Selesaikan dhikr pagi dan petang",
		Icon:        "BicepsFlexed",
		Unlocked:    hasMorning && hasEvening,
	})

	return achievements, nil
}

func (s *appService) GetWeeklyActivity(ctx context.Context, userID string) ([]serviceinterface.WeeklyActivityItem, error) {
	now := time.Now()
	activity := make([]serviceinterface.WeeklyActivityItem, 7)

	// We'll approximate active dates by grabbing progress, quiz, and dhikr history
	// In a real app we'd query an activity_log table.
	// We will just do a simple aggregation: if there's any progress last accessed on the day, it's a hit.

	progress, _ := s.repo.GetUserProgress(ctx, userID, nil)
	quizzes, _ := s.repo.GetUserQuizAttempts(ctx, userID, nil)

	activeDates := make(map[string]bool)
	for _, p := range progress {
		activeDates[p.LastAccessed.Format("2006-01-02")] = true
	}
	for _, q := range quizzes {
		activeDates[q.CompletedAt.Format("2006-01-02")] = true
	}

	// Calculate last 7 days ending today
	dayNames := []string{"Min", "Sen", "Sel", "Rab", "Kam", "Jum", "Sab"}

	for i := 6; i >= 0; i-- {
		targetDate := now.AddDate(0, 0, i-6)
		dateStr := targetDate.Format("2006-01-02")

		// Check Dhikr for that date
		hasDhikr := false
		if dhikrCounters, err := s.repo.GetDhikrCountersForDate(ctx, userID, dateStr); err == nil && len(dhikrCounters) > 0 {
			for _, dc := range dhikrCounters {
				if dc.Completed {
					hasDhikr = true
					break
				}
			}
		}

		isActive := activeDates[dateStr] || hasDhikr
		activity[i] = serviceinterface.WeeklyActivityItem{
			Day:       dayNames[targetDate.Weekday()],
			Completed: isActive,
		}
	}

	return activity, nil
}
