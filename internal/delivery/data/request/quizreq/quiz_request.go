package quizreq

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
