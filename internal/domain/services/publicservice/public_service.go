package publicservice

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
)

const publicUserAgent = "Hijaiyah-App/1.0"

type publicService struct {
	httpClient *http.Client
}

type reverseGeocodeResult struct {
	City     *string
	Country  *string
	District *string
	Label    *string
}

type aladhanResponse struct {
	Code int `json:"code"`
	Data struct {
		Timings struct {
			Fajr    string `json:"Fajr"`
			Sunrise string `json:"Sunrise"`
			Dhuhr   string `json:"Dhuhr"`
			Asr     string `json:"Asr"`
			Maghrib string `json:"Maghrib"`
			Isha    string `json:"Isha"`
		} `json:"timings"`
		Date struct {
			Gregorian struct {
				Date string `json:"date"`
			} `json:"gregorian"`
		} `json:"date"`
		Meta struct {
			Timezone string `json:"timezone"`
			Method   struct {
				Name string `json:"name"`
			} `json:"method"`
		} `json:"meta"`
	} `json:"data"`
}

type nominatimResponse struct {
	DisplayName string `json:"display_name"`
	Address     struct {
		City          string `json:"city"`
		Town          string `json:"town"`
		Village       string `json:"village"`
		Suburb        string `json:"suburb"`
		County        string `json:"county"`
		StateDistrict string `json:"state_district"`
		Country       string `json:"country"`
	} `json:"address"`
}

func NewPublicService(httpClient *http.Client) serviceinterface.PublicService {
	client := httpClient
	if client == nil {
		client = &http.Client{Timeout: 12 * time.Second}
	}
	return &publicService{httpClient: client}
}

func (s *publicService) GetPrayerTimes(
	ctx context.Context,
	latitude float64,
	longitude float64,
	dateValue string,
	method string,
) (serviceinterface.PrayerTimes, serviceinterface.PrayerTimesMeta) {
	times, meta, err := s.fetchPrayerTimes(ctx, latitude, longitude, dateValue, method)
	if err == nil {
		return times, meta
	}

	return calculateFallbackPrayerTimes(latitude, longitude, dateValue), serviceinterface.PrayerTimesMeta{
		Method: "Fallback calculation",
		Note:   "External API unavailable, using approximate times",
	}
}

func (s *publicService) fetchPrayerTimes(
	ctx context.Context,
	latitude float64,
	longitude float64,
	dateValue string,
	method string,
) (serviceinterface.PrayerTimes, serviceinterface.PrayerTimesMeta, error) {
	apiURL := fmt.Sprintf("https://api.aladhan.com/v1/timings/%s", dateValue)
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	params.Set("method", method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?%s", apiURL, params.Encode()), nil)
	if err != nil {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, err
	}
	req.Header.Set("User-Agent", publicUserAgent)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, fmt.Errorf("aladhan responded with status %d", res.StatusCode)
	}

	var payload aladhanResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, err
	}
	if payload.Code != http.StatusOK {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, fmt.Errorf("aladhan response code %d", payload.Code)
	}

	place := s.reverseGeocode(ctx, latitude, longitude)
	city := place.City
	if city == nil {
		timezone := strings.TrimSpace(payload.Data.Meta.Timezone)
		if timezone != "" {
			city = &timezone
		}
	}

	return serviceinterface.PrayerTimes{
			Fajr:    payload.Data.Timings.Fajr,
			Sunrise: payload.Data.Timings.Sunrise,
			Dhuhr:   payload.Data.Timings.Dhuhr,
			Asr:     payload.Data.Timings.Asr,
			Maghrib: payload.Data.Timings.Maghrib,
			Isha:    payload.Data.Timings.Isha,
			Date:    payload.Data.Date.Gregorian.Date,
			Location: serviceinterface.PrayerLocation{
				Latitude:  latitude,
				Longitude: longitude,
				City:      city,
				Country:   place.Country,
				District:  place.District,
				Label:     place.Label,
			},
		}, serviceinterface.PrayerTimesMeta{
			Method:            payload.Data.Meta.Method.Name,
			Timezone:          payload.Data.Meta.Timezone,
			CalculationMethod: method,
		}, nil
}

func (s *publicService) reverseGeocode(ctx context.Context, latitude float64, longitude float64) reverseGeocodeResult {
	reverseURL, err := url.Parse("https://nominatim.openstreetmap.org/reverse")
	if err != nil {
		return reverseGeocodeResult{}
	}

	query := reverseURL.Query()
	query.Set("format", "jsonv2")
	query.Set("lat", strconv.FormatFloat(latitude, 'f', -1, 64))
	query.Set("lon", strconv.FormatFloat(longitude, 'f', -1, 64))
	query.Set("zoom", "10")
	query.Set("addressdetails", "1")
	reverseURL.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reverseURL.String(), nil)
	if err != nil {
		return reverseGeocodeResult{}
	}
	req.Header.Set("User-Agent", publicUserAgent)

	res, err := s.httpClient.Do(req)
	if err != nil {
		return reverseGeocodeResult{}
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return reverseGeocodeResult{}
	}

	var payload nominatimResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return reverseGeocodeResult{}
	}

	city := firstNonEmpty(payload.Address.City, payload.Address.Town, payload.Address.Village)
	district := firstNonEmpty(payload.Address.Suburb, payload.Address.County, payload.Address.StateDistrict)
	country := strings.TrimSpace(payload.Address.Country)
	label := strings.TrimSpace(payload.DisplayName)

	return reverseGeocodeResult{
		City:     strPtr(city),
		Country:  strPtr(country),
		District: strPtr(district),
		Label:    strPtr(label),
	}
}

func calculateFallbackPrayerTimes(latitude float64, longitude float64, dateValue string) serviceinterface.PrayerTimes {
	solarNoon := 12 - (longitude / 15)

	return serviceinterface.PrayerTimes{
		Fajr:    formatTime(solarNoon - 1.5),
		Sunrise: formatTime(solarNoon - 1),
		Dhuhr:   formatTime(solarNoon),
		Asr:     formatTime(solarNoon + 3),
		Maghrib: formatTime(solarNoon + 6),
		Isha:    formatTime(solarNoon + 7.5),
		Date:    dateValue,
		Location: serviceinterface.PrayerLocation{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}
}

func formatTime(decimalHours float64) string {
	hours := int(math.Floor(decimalHours))
	minutes := int(math.Floor((decimalHours - float64(hours)) * 60))
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func strPtr(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
