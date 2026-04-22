package habitreq

type CreateHabitRequest struct {
	Name            string `json:"name"`
	Category        string `json:"category"`
	ReminderTime    string `json:"reminderTime"`
	ReminderEnabled *bool  `json:"reminderEnabled,omitempty"`
}

type UpdateHabitRequest struct {
	Name            *string `json:"name,omitempty"`
	Category        *string `json:"category,omitempty"`
	ReminderTime    *string `json:"reminderTime,omitempty"`
	ReminderEnabled *bool   `json:"reminderEnabled,omitempty"`
}

type SetHabitCompletionRequest struct {
	HabitID   string `json:"habitId"`
	Date      string `json:"date"`
	Completed *bool  `json:"completed,omitempty"`
}
