package progressreq

type UpsertProgressRequest struct {
	Module    string `json:"module"`
	ItemID    string `json:"itemId"`
	Progress  int    `json:"progress"`
	Completed *bool  `json:"completed,omitempty"`
	Score     *int   `json:"score,omitempty"`
	TimeSpent *int   `json:"timeSpent,omitempty"`
}
