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
	"sync"
	"time"

	"github.com/awahids/bn-server/internal/domain/services/serviceinterface"
)

const publicUserAgent = "Hijaiyah-App/1.0"

const (
	// Prayer times only change per day and location; caching them avoids an
	// Aladhan round trip on every home-page load.
	prayerTimesCacheTTL    = 6 * time.Hour
	prayerFallbackCacheTTL = 10 * time.Minute
	// Nominatim's usage policy allows ~1 request/second, so place names are
	// cached aggressively (they practically never change).
	geocodeCacheTTL      = 24 * time.Hour
	geocodeEmptyCacheTTL = 10 * time.Minute
	maxCacheEntries      = 5000
)

type publicService struct {
	httpClient *http.Client

	prayerCache  *ttlCache[cachedPrayerTimes]
	geocodeCache *ttlCache[reverseGeocodeResult]
}

type cachedPrayerTimes struct {
	times serviceinterface.PrayerTimes
	meta  serviceinterface.PrayerTimesMeta
}

type ttlCacheEntry[T any] struct {
	value     T
	expiresAt time.Time
}

// ttlCache is a small, size-capped, concurrency-safe in-memory cache.
type ttlCache[T any] struct {
	mu      sync.Mutex
	entries map[string]ttlCacheEntry[T]
}

func newTTLCache[T any]() *ttlCache[T] {
	return &ttlCache[T]{entries: make(map[string]ttlCacheEntry[T])}
}

func (c *ttlCache[T]) get(key string, now time.Time) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.entries[key]
	if !ok || now.After(entry.expiresAt) {
		var zero T
		return zero, false
	}
	return entry.value, true
}

func (c *ttlCache[T]) set(key string, value T, ttl time.Duration, now time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if len(c.entries) >= maxCacheEntries {
		for k, entry := range c.entries {
			if now.After(entry.expiresAt) {
				delete(c.entries, k)
			}
		}
		if len(c.entries) >= maxCacheEntries {
			c.entries = make(map[string]ttlCacheEntry[T])
		}
	}
	c.entries[key] = ttlCacheEntry[T]{value: value, expiresAt: now.Add(ttl)}
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
	return &publicService{
		httpClient:   client,
		prayerCache:  newTTLCache[cachedPrayerTimes](),
		geocodeCache: newTTLCache[reverseGeocodeResult](),
	}
}

func (s *publicService) GetPrayerTimes(
	ctx context.Context,
	latitude float64,
	longitude float64,
	dateValue string,
	method string,
) (serviceinterface.PrayerTimes, serviceinterface.PrayerTimesMeta) {
	// ~110 m precision: nearby users share the same entry.
	cacheKey := fmt.Sprintf("%s|%s|%.3f|%.3f", dateValue, method, latitude, longitude)
	if cached, ok := s.prayerCache.get(cacheKey, time.Now()); ok {
		cached.times.Location.Latitude = latitude
		cached.times.Location.Longitude = longitude
		return cached.times, cached.meta
	}

	times, meta, err := s.fetchPrayerTimes(ctx, latitude, longitude, dateValue, method)
	if err == nil {
		s.prayerCache.set(cacheKey, cachedPrayerTimes{times: times, meta: meta}, prayerTimesCacheTTL, time.Now())
		return times, meta
	}

	fallbackTimes := calculateFallbackPrayerTimes(latitude, longitude, dateValue)
	fallbackMeta := serviceinterface.PrayerTimesMeta{
		Method: "Fallback calculation",
		Note:   "External API unavailable, using approximate times",
	}
	// Cache briefly so an Aladhan outage doesn't add a failing round trip to
	// every request, while still retrying the provider soon.
	s.prayerCache.set(cacheKey, cachedPrayerTimes{times: fallbackTimes, meta: fallbackMeta}, prayerFallbackCacheTTL, time.Now())
	return fallbackTimes, fallbackMeta
}

func (s *publicService) fetchPrayerTimes(
	ctx context.Context,
	latitude float64,
	longitude float64,
	dateValue string,
	method string,
) (serviceinterface.PrayerTimes, serviceinterface.PrayerTimesMeta, error) {
	// Aladhan expects the path date as DD-MM-YYYY.
	providerDate := dateValue
	if parsed, err := time.Parse("2006-01-02", dateValue); err == nil {
		providerDate = parsed.Format("02-01-2006")
	}
	apiURL := fmt.Sprintf("https://api.aladhan.com/v1/timings/%s", providerDate)
	params := url.Values{}
	params.Set("latitude", strconv.FormatFloat(latitude, 'f', -1, 64))
	params.Set("longitude", strconv.FormatFloat(longitude, 'f', -1, 64))
	params.Set("method", method)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s?%s", apiURL, params.Encode()), nil)
	if err != nil {
		return serviceinterface.PrayerTimes{}, serviceinterface.PrayerTimesMeta{}, err
	}
	req.Header.Set("User-Agent", publicUserAgent)

	// Look up the place name concurrently instead of after Aladhan answers.
	placeCh := make(chan reverseGeocodeResult, 1)
	go func() {
		placeCh <- s.reverseGeocode(ctx, latitude, longitude)
	}()

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

	place := <-placeCh
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
			// Keep the requested YYYY-MM-DD format (Aladhan returns DD-MM-YYYY) so
			// the response matches the fallback and the endpoint's documented format.
			Date: dateValue,
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
	// ~1 km precision is plenty for a city/district label.
	cacheKey := fmt.Sprintf("%.2f|%.2f", latitude, longitude)
	if cached, ok := s.geocodeCache.get(cacheKey, time.Now()); ok {
		return cached
	}

	place := s.fetchReverseGeocode(ctx, latitude, longitude)
	ttl := geocodeCacheTTL
	if place.City == nil && place.District == nil && place.Label == nil {
		// Failed or empty lookups are retried sooner.
		ttl = geocodeEmptyCacheTTL
	}
	s.geocodeCache.set(cacheKey, place, ttl, time.Now())
	return place
}

func (s *publicService) fetchReverseGeocode(ctx context.Context, latitude float64, longitude float64) reverseGeocodeResult {
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

// calculateFallbackPrayerTimes approximates prayer times astronomically when the
// external provider is unavailable. Times are expressed in the approximate local
// timezone derived from longitude (e.g. WIB for western Indonesia), using the
// Kemenag RI convention: Fajr 20°, Isha 18°, Shafi'i Asr.
func calculateFallbackPrayerTimes(latitude float64, longitude float64, dateValue string) serviceinterface.PrayerTimes {
	date, err := time.Parse("2006-01-02", dateValue)
	if err != nil {
		date = time.Now().UTC()
	}

	const degToRad = math.Pi / 180
	dayOfYear := float64(date.YearDay())
	gamma := 2 * math.Pi / 365 * (dayOfYear - 1)
	equationOfTime := 229.18 * (0.000075 + 0.001868*math.Cos(gamma) - 0.032077*math.Sin(gamma) -
		0.014615*math.Cos(2*gamma) - 0.040849*math.Sin(2*gamma))
	declination := 0.006918 - 0.399912*math.Cos(gamma) + 0.070257*math.Sin(gamma) -
		0.006758*math.Cos(2*gamma) + 0.000907*math.Sin(2*gamma) -
		0.002697*math.Cos(3*gamma) + 0.00148*math.Sin(3*gamma)

	timezoneOffset := math.Round(longitude / 15)
	solarNoon := 12 + timezoneOffset - longitude/15 - equationOfTime/60
	latRad := latitude * degToRad

	// hourAngle returns the hours between solar noon and the moment the sun
	// reaches the given altitude (degrees, negative below the horizon).
	hourAngle := func(altitude float64) float64 {
		cosH := (math.Sin(altitude*degToRad) - math.Sin(latRad)*math.Sin(declination)) /
			(math.Cos(latRad) * math.Cos(declination))
		cosH = math.Max(-1, math.Min(1, cosH))
		return math.Acos(cosH) / degToRad / 15
	}

	asrAltitude := math.Atan(1/(1+math.Tan(math.Abs(latRad-declination)))) / degToRad
	// Kemenag adds a few minutes of ihtiyat (safety margin) to each time.
	const ihtiyat = 2.0 / 60

	return serviceinterface.PrayerTimes{
		Fajr:    formatTime(solarNoon - hourAngle(-20) + ihtiyat),
		Sunrise: formatTime(solarNoon - hourAngle(-0.833) - ihtiyat),
		Dhuhr:   formatTime(solarNoon + ihtiyat),
		Asr:     formatTime(solarNoon + hourAngle(asrAltitude) + ihtiyat),
		Maghrib: formatTime(solarNoon + hourAngle(-0.833) + ihtiyat),
		Isha:    formatTime(solarNoon + hourAngle(-18) + ihtiyat),
		Date:    dateValue,
		Location: serviceinterface.PrayerLocation{
			Latitude:  latitude,
			Longitude: longitude,
		},
	}
}

func formatTime(decimalHours float64) string {
	totalMinutes := int(math.Round(decimalHours * 60))
	totalMinutes = ((totalMinutes % 1440) + 1440) % 1440
	return fmt.Sprintf("%02d:%02d", totalMinutes/60, totalMinutes%60)
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
