package serviceinterface

import "context"

type PrayerTimes struct {
	Fajr     string         `json:"fajr"`
	Sunrise  string         `json:"sunrise"`
	Dhuhr    string         `json:"dhuhr"`
	Asr      string         `json:"asr"`
	Maghrib  string         `json:"maghrib"`
	Isha     string         `json:"isha"`
	Date     string         `json:"date"`
	Location PrayerLocation `json:"location"`
}

type PrayerLocation struct {
	City      *string `json:"city,omitempty"`
	Country   *string `json:"country,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	District  *string `json:"district,omitempty"`
	Label     *string `json:"label,omitempty"`
}

type PrayerTimesMeta struct {
	Method            string `json:"method"`
	Timezone          string `json:"timezone,omitempty"`
	CalculationMethod string `json:"calculationMethod,omitempty"`
	Note              string `json:"note,omitempty"`
}

type PublicService interface {
	GetPrayerTimes(
		ctx context.Context,
		latitude float64,
		longitude float64,
		dateValue string,
		method string,
	) (PrayerTimes, PrayerTimesMeta)
}
