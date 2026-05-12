package handlerutil

import "testing"

func TestIsValidDate(t *testing.T) {
	testCases := []struct {
		name  string
		value string
		valid bool
	}{
		{name: "valid regular date", value: "2026-04-23", valid: true},
		{name: "valid leap day", value: "2024-02-29", valid: true},
		{name: "invalid leap day", value: "2026-02-29", valid: false},
		{name: "invalid day", value: "2026-04-31", valid: false},
		{name: "invalid month", value: "2026-13-01", valid: false},
		{name: "invalid format no padding", value: "2026-4-1", valid: false},
		{name: "invalid format text", value: "not-a-date", valid: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidDate(tc.value); got != tc.valid {
				t.Fatalf("expected %v, got %v for value %q", tc.valid, got, tc.value)
			}
		})
	}
}

func TestIsValidProgressModule(t *testing.T) {
	testCases := []struct {
		name   string
		module string
		valid  bool
	}{
		{name: "hijaiyah", module: "hijaiyah", valid: true},
		{name: "tajwid", module: "tajwid", valid: true},
		{name: "quran", module: "quran", valid: true},
		{name: "dhikr", module: "dhikr", valid: true},
		{name: "quiz", module: "quiz", valid: true},
		{name: "hafalan", module: "hafalan", valid: true},
		{name: "invalid", module: "unknown", valid: false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsValidProgressModule(tc.module); got != tc.valid {
				t.Fatalf("expected %v, got %v for module %q", tc.valid, got, tc.module)
			}
		})
	}
}
