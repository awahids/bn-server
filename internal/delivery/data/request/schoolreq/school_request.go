package schoolreq

type CreateSchoolRequest struct {
	Name          string  `json:"name"`
	Location      string  `json:"location"`
	Jenjang       string  `json:"jenjang"`
	StatusSekolah string  `json:"statusSekolah"`
	MonthlyFee    *int    `json:"monthlyFee"`
	MapURL        string  `json:"mapUrl"`
	Contact       *string `json:"contact,omitempty"`
	Description   *string `json:"description,omitempty"`
}
