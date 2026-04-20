package dhikrreq

type UpsertDhikrCounterRequest struct {
	DhikrID   string `json:"dhikrId"`
	Count     int    `json:"count"`
	Target    *int   `json:"target,omitempty"`
	Date      string `json:"date"`
	Session   string `json:"session"`
	Completed *bool  `json:"completed,omitempty"`
}
