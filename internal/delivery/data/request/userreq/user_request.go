package userreq

type UpdateUserRequest struct {
	Name          *string        `json:"name"`
	Username      *string        `json:"username"`
	Streak        *int           `json:"streak"`
	DailyProgress *int           `json:"dailyProgress"`
	Preferences   map[string]any `json:"preferences"`
}
