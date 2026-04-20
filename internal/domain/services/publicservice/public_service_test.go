package publicservice

import (
	"bytes"
	"context"
	"errors"
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
}

func httpResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}
