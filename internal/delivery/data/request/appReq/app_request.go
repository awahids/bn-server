package appreq

type UpdateUserRequest struct {
	Name          *string        `json:"name"`
	Username      *string        `json:"username"`
	Streak        *int           `json:"streak"`
	DailyProgress *int           `json:"dailyProgress"`
	Preferences   map[string]any `json:"preferences"`
}

type UpsertProgressRequest struct {
	Module    string `json:"module"`
	ItemID    string `json:"itemId"`
	Progress  int    `json:"progress"`
	Completed *bool  `json:"completed,omitempty"`
	Score     *int   `json:"score,omitempty"`
	TimeSpent *int   `json:"timeSpent,omitempty"`
}

type CreateBookmarkRequest struct {
	Type      string  `json:"type"`
	ContentID string  `json:"contentId"`
	Note      *string `json:"note"`
}

type UpsertDhikrCounterRequest struct {
	DhikrID   string `json:"dhikrId"`
	Count     int    `json:"count"`
	Target    *int   `json:"target,omitempty"`
	Date      string `json:"date"`
	Session   string `json:"session"`
	Completed *bool  `json:"completed,omitempty"`
}

type QuizAnswerRequest struct {
	QuestionID    string `json:"questionId"`
	UserAnswer    string `json:"userAnswer"`
	CorrectAnswer string `json:"correctAnswer"`
	IsCorrect     bool   `json:"isCorrect"`
	TimeSpent     *int   `json:"timeSpent,omitempty"`
}

type CreateQuizAttemptRequest struct {
	Category       string              `json:"category"`
	Score          int                 `json:"score"`
	TotalQuestions int                 `json:"totalQuestions"`
	TimeSpent      int                 `json:"timeSpent"`
	Answers        []QuizAnswerRequest `json:"answers"`
}
