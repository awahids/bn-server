package publicservice

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestGetPrayerTimes_Success(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.Host {
			case "api.aladhan.com":
				if req.URL.Path != "/v1/timings/20-04-2026" {
					return nil, fmt.Errorf("unexpected aladhan path %s", req.URL.Path)
				}
				body := `{"code":200,"data":{"timings":{"Fajr":"04:30","Sunrise":"05:45","Dhuhr":"12:05","Asr":"15:25","Maghrib":"18:05","Isha":"19:15"},"date":{"gregorian":{"date":"20-04-2026"}},"meta":{"timezone":"Asia/Jakarta","method":{"name":"Muslim World League"}}}}`
				return httpResponse(http.StatusOK, body), nil
			case "nominatim.openstreetmap.org":
				body := `{"display_name":"Jakarta, Indonesia","address":{"city":"Jakarta","country":"Indonesia","suburb":"Menteng"}}`
				return httpResponse(http.StatusOK, body), nil
			default:
				return nil, errors.New("unexpected host")
			}
		}),
	}

	svc := NewPublicService(client)
	times, meta := svc.GetPrayerTimes(context.Background(), -6.2, 106.8, "2026-04-20", "2")

	if times.Date != "2026-04-20" {
		t.Fatalf("expected date 2026-04-20, got %s", times.Date)
	}
	if times.Fajr != "04:30" {
		t.Fatalf("expected fajr 04:30, got %s", times.Fajr)
	}
	if times.Location.City == nil || *times.Location.City != "Jakarta" {
		t.Fatalf("expected city Jakarta, got %+v", times.Location.City)
	}
	if meta.Method != "Muslim World League" {
		t.Fatalf("expected method Muslim World League, got %s", meta.Method)
	}
	if meta.Note != "" {
		t.Fatalf("expected empty note on success, got %s", meta.Note)
	}
}

func TestGetPrayerTimes_CachesProviderResult(t *testing.T) {
	var aladhanCalls, nominatimCalls int
	client := &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			switch req.URL.Host {
			case "api.aladhan.com":
				aladhanCalls++
				body := `{"code":200,"data":{"timings":{"Fajr":"04:30","Sunrise":"05:45","Dhuhr":"12:05","Asr":"15:25","Maghrib":"18:05","Isha":"19:15"},"date":{"gregorian":{"date":"20-04-2026"}},"meta":{"timezone":"Asia/Jakarta","method":{"name":"Kemenag"}}}}`
				return httpResponse(http.StatusOK, body), nil
			case "nominatim.openstreetmap.org":
				nominatimCalls++
				return httpResponse(http.StatusOK, `{"display_name":"Jakarta","address":{"city":"Jakarta"}}`), nil
			default:
				return nil, errors.New("unexpected host")
			}
		}),
	}

	svc := NewPublicService(client)
	first, _ := svc.GetPrayerTimes(context.Background(), -6.2001, 106.8001, "2026-04-20", "20")
	second, _ := svc.GetPrayerTimes(context.Background(), -6.2002, 106.8002, "2026-04-20", "20")

	if aladhanCalls != 1 || nominatimCalls != 1 {
		t.Fatalf("expected one provider call each, got aladhan=%d nominatim=%d", aladhanCalls, nominatimCalls)
	}
	if second.Fajr != first.Fajr {
		t.Fatalf("expected cached times, got %s vs %s", second.Fajr, first.Fajr)
	}
	if second.Location.Latitude != -6.2002 {
		t.Fatalf("expected cached response to echo the requested latitude, got %v", second.Location.Latitude)
	}
}

func TestGetPrayerTimes_FallbackWhenProviderFails(t *testing.T) {
	client := &http.Client{
		Timeout: time.Second,
		Transport: roundTripFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.Host == "api.aladhan.com" {
				return nil, errors.New("provider unavailable")
			}
			return nil, errors.New("unexpected request")
		}),
	}

	svc := NewPublicService(client)
	times, meta := svc.GetPrayerTimes(context.Background(), -6.2, 106.8, "2026-04-20", "2")

	if times.Date != "2026-04-20" {
		t.Fatalf("expected fallback date 2026-04-20, got %s", times.Date)
	}
	if meta.Method != "Fallback calculation" {
		t.Fatalf("expected fallback method, got %s", meta.Method)
	}
	if meta.Note == "" {
		t.Fatal("expected fallback note to be present")
	}

	// Jakarta on 2026-04-20 (WIB): Subuh ~04:3x, Dzuhur ~11:5x, Maghrib ~17:5x.
	expectations := map[string]struct {
		value string
		from  string
		to    string
	}{
		"fajr":    {times.Fajr, "04:20", "04:45"},
		"dhuhr":   {times.Dhuhr, "11:45", "12:00"},
		"asr":     {times.Asr, "15:00", "15:25"},
		"maghrib": {times.Maghrib, "17:45", "18:05"},
		"isha":    {times.Isha, "18:55", "19:15"},
	}
	for name, exp := range expectations {
		if exp.value < exp.from || exp.value > exp.to {
			t.Fatalf("expected fallback %s between %s and %s, got %s", name, exp.from, exp.to, exp.value)
		}
	}
}

func httpResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func TestCalculateFallbackPrayerTimes_MatchesPublishedJakartaSchedule(t *testing.T) {
	// Published Jakarta schedule for 2026-10-29 (6°10' S, 106°49' E):
	// Subuh 04:10, Terbit 05:25, Dzuhur 11:39, Ashar 14:53, Maghrib 17:48, Isya 18:59.
	times := calculateFallbackPrayerTimes(-6.1667, 106.8167, "2026-10-29")
	expected := map[string][2]string{
		"Fajr":    {times.Fajr, "04:10"},
		"Sunrise": {times.Sunrise, "05:25"},
		"Dhuhr":   {times.Dhuhr, "11:39"},
		"Asr":     {times.Asr, "14:53"},
		"Maghrib": {times.Maghrib, "17:48"},
		"Isha":    {times.Isha, "18:59"},
	}
	toMinutes := func(value string) int {
		parsed, err := time.Parse("15:04", value)
		if err != nil {
			t.Fatalf("invalid time %q: %v", value, err)
		}
		return parsed.Hour()*60 + parsed.Minute()
	}
	for name, pair := range expected {
		diff := toMinutes(pair[0]) - toMinutes(pair[1])
		if diff < -2 || diff > 2 {
			t.Errorf("%s = %s, want %s ±2 min", name, pair[0], pair[1])
		}
	}
}
