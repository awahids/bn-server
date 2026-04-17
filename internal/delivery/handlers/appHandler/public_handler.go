package apphandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"bn-mobile/server/internal/delivery/data/response"

	"github.com/gin-gonic/gin"
)

const (
	publicRequestTimeout = 15 * time.Second
	audioUserAgent       = "Hijaiyah-App/1.0"
)

var allowedAudioDomains = []string{"localhost", "127.0.0.1"}

type PublicHandler struct {
	httpClient *http.Client
}

type PrayerTimes struct {
	Fajr     string         `json:"fajr"`
	Sunrise  string         `json:"sunrise"`
	Dhuhr    string         `json:"dhuhr"`
	Asr      string         `json:"asr"`
	Maghrib  string         `json:"maghrib"`
	Isha     string         `json:"isha"`
	Date     string         `json:"date"`
	Location prayerLocation `json:"location"`
}

type prayerLocation struct {
	City      *string `json:"city,omitempty"`
	Country   *string `json:"country,omitempty"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	District  *string `json:"district,omitempty"`
	Label     *string `json:"label,omitempty"`
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

func NewPublicHandler() *PublicHandler {
	return &PublicHandler{
		httpClient: &http.Client{
			Timeout: 12 * time.Second,
		},
	}
}

// GetAudioProxy godoc
// @Summary Audio proxy
// @Description Proxy local/allowed audio URL for client playback.
// @Tags Public
// @Produce octet-stream
// @Param url query string true "Audio source URL or local path"
// @Success 200 {file} file
// @Failure 400 {object} response.APIResponse
// @Failure 403 {object} response.APIResponse
// @Failure 500 {object} response.APIResponse
// @Router /audio-proxy [get]
func (h *PublicHandler) GetAudioProxy(c *gin.Context) {
	audioURL := strings.TrimSpace(c.Query("url"))
	if audioURL == "" {
		response.Failed(c, http.StatusBadRequest, "audio URL is required", "audio URL is required")
		return
	}

	isLocalFile := strings.HasPrefix(audioURL, "/") || strings.HasPrefix(audioURL, "./") || strings.HasPrefix(audioURL, "../")
	if isLocalFile {
		filePath := audioURL
		if !strings.HasPrefix(filePath, "/") {
			filePath = "/" + filePath
		}
		if strings.Contains(filePath, "..") || strings.Contains(filePath, "~") {
			response.Failed(c, http.StatusBadRequest, "invalid file path", "invalid file path")
			return
		}
		c.Redirect(http.StatusTemporaryRedirect, filePath)
		return
	}

	parsedURL, err := url.Parse(audioURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		response.Failed(c, http.StatusBadRequest, "invalid URL format", "invalid URL format")
		return
	}

	domainAllowed := false
	for _, domain := range allowedAudioDomains {
		if strings.Contains(parsedURL.Hostname(), domain) {
			domainAllowed = true
			break
		}
	}
	if !domainAllowed {
		response.Failed(c, http.StatusForbidden, "domain not allowed", "domain not allowed")
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), publicRequestTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, audioURL, nil)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid URL format", err.Error())
		return
	}
	req.Header.Set("User-Agent", audioUserAgent)

	res, err := h.httpClient.Do(req)
	if err != nil {
		response.Failed(c, http.StatusInternalServerError, "failed to fetch audio file", err.Error())
		return
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		response.Failed(c, res.StatusCode, "failed to fetch audio file", "failed to fetch audio file")
		return
	}

	contentType := strings.TrimSpace(res.Header.Get("Content-Type"))
	if contentType == "" {
		contentType = "audio/mpeg"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "Content-Type")
	c.Status(http.StatusOK)

	if _, err := io.Copy(c.Writer, res.Body); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}

// GetPrayerTimes godoc
// @Summary Get prayer times
// @Description Get prayer times by coordinates and optional date/method.
// @Tags Public
// @Produce json
// @Param lat query number true "Latitude"
// @Param lng query number true "Longitude"
// @Param date query string false "Date in YYYY-MM-DD (default today)"
// @Param method query string false "Calculation method (default 2)"
// @Success 200 {object} map[string]any
// @Failure 400 {object} response.APIResponse
// @Router /prayer-times [get]
func (h *PublicHandler) GetPrayerTimes(c *gin.Context) {
	latParam := strings.TrimSpace(c.Query("lat"))
	lngParam := strings.TrimSpace(c.Query("lng"))
	if latParam == "" || lngParam == "" {
		response.Failed(c, http.StatusBadRequest, "latitude and longitude are required", "latitude and longitude are required")
		return
	}

	latitude, err := strconv.ParseFloat(latParam, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "invalid latitude")
		return
	}
	longitude, err := strconv.ParseFloat(lngParam, 64)
	if err != nil {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "invalid longitude")
		return
	}

	if latitude < -90 || latitude > 90 || longitude < -180 || longitude > 180 {
		response.Failed(c, http.StatusBadRequest, "invalid latitude or longitude", "latitude must be between -90 and 90, longitude between -180 and 180")
		return
	}

	dateValue := strings.TrimSpace(c.Query("date"))
	if dateValue == "" {
		dateValue = time.Now().UTC().Format("2006-01-02")
	}
	if !dateRegex.MatchString(dateValue) {
		response.Failed(c, http.StatusBadRequest, "invalid date format", "date must be in YYYY-MM-DD format")
		return
	}

	method := strings.TrimSpace(c.Query("method"))
	if method == "" {
		method = "2"
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), publicRequestTimeout)
	defer cancel()

	prayerTimes, meta, err := h.fetchPrayerTimes(ctx, latitude, longitude, dateValue, method)
	if err != nil {
		fallback := calculateFallbackPrayerTimes(latitude, longitude, dateValue)
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    fallback,
			"meta": gin.H{
				"method": "Fallback calculation",
				"note":   "External API unavailable, using approximate times",
			},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    prayerTimes,
		"meta":    meta,
	})
}

func (h *PublicHandler) fetchPrayerTimes(
	ctx context.Context,
	latitude float64,
	longitude float64,
	dateValue string,
	method string,
) (PrayerTimes, gin.H, error) {
	apiURL := fmt.Sprintf("http://api.aladhan.com/v1/timings/%s", dateValue)
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	params.Set("method", method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?%s", apiURL, params.Encode()), nil)
	if err != nil {
		return PrayerTimes{}, nil, err
	}
	req.Header.Set("User-Agent", audioUserAgent)

	res, err := h.httpClient.Do(req)
	if err != nil {
		return PrayerTimes{}, nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusMultipleChoices {
		return PrayerTimes{}, nil, fmt.Errorf("aladhan responded with status %d", res.StatusCode)
	}

	var payload aladhanResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		return PrayerTimes{}, nil, err
	}
	if payload.Code != http.StatusOK {
		return PrayerTimes{}, nil, fmt.Errorf("aladhan response code %d", payload.Code)
	}

	place := h.reverseGeocode(ctx, latitude, longitude)

	city := place.City
	if city == nil {
		timezone := strings.TrimSpace(payload.Data.Meta.Timezone)
		if timezone != "" {
			city = &timezone
		}
	}

	result := PrayerTimes{
		Fajr:    payload.Data.Timings.Fajr,
		Sunrise: payload.Data.Timings.Sunrise,
		Dhuhr:   payload.Data.Timings.Dhuhr,
		Asr:     payload.Data.Timings.Asr,
		Maghrib: payload.Data.Timings.Maghrib,
		Isha:    payload.Data.Timings.Isha,
		Date:    payload.Data.Date.Gregorian.Date,
		Location: prayerLocation{
			Latitude:  latitude,
			Longitude: longitude,
			City:      city,
			Country:   place.Country,
			District:  place.District,
			Label:     place.Label,
		},
	}

	meta := gin.H{
		"method":            payload.Data.Meta.Method.Name,
		"timezone":          payload.Data.Meta.Timezone,
		"calculationMethod": method,
	}
	return result, meta, nil
}

func (h *PublicHandler) reverseGeocode(ctx context.Context, latitude float64, longitude float64) reverseGeocodeResult {
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
	req.Header.Set("User-Agent", audioUserAgent)

	res, err := h.httpClient.Do(req)
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

func calculateFallbackPrayerTimes(latitude float64, longitude float64, dateValue string) PrayerTimes {
	solarNoon := 12 - (longitude / 15)

	return PrayerTimes{
		Fajr:    formatTime(solarNoon - 1.5),
		Sunrise: formatTime(solarNoon - 1),
		Dhuhr:   formatTime(solarNoon),
		Asr:     formatTime(solarNoon + 3),
		Maghrib: formatTime(solarNoon + 6),
		Isha:    formatTime(solarNoon + 7.5),
		Date:    dateValue,
		Location: prayerLocation{
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
