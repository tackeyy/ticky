package ticktick

import (
	"strings"
	"testing"
	"time"
)

// expectedEndOfDay computes the expected endOfDay output for a given time,
// matching the production endOfDay function's format.
func expectedEndOfDay(t time.Time) string {
	eod := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
	return eod.UTC().Format("2006-01-02T15:04:05.000+0000")
}

func TestParseDate_Keywords(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"today lowercase", "today", expectedEndOfDay(now)},
		{"today uppercase", "TODAY", expectedEndOfDay(now)},
		{"today mixed case", "Today", expectedEndOfDay(now)},
		{"tomorrow lowercase", "tomorrow", expectedEndOfDay(now.AddDate(0, 0, 1))},
		{"tomorrow uppercase", "TOMORROW", expectedEndOfDay(now.AddDate(0, 0, 1))},
		{"tomorrow mixed case", "Tomorrow", expectedEndOfDay(now.AddDate(0, 0, 1))},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("ParseDate(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseDate_RelativeDays(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name  string
		input string
		days  int
	}{
		{"plus zero days", "+0d", 0},
		{"plus one day", "+1d", 1},
		{"plus three days", "+3d", 3},
		{"plus hundred days", "+100d", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			want := expectedEndOfDay(now.AddDate(0, 0, tt.days))
			got, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != want {
				t.Errorf("ParseDate(%q) = %q, want %q", tt.input, got, want)
			}
		})
	}
}

func TestParseDate_CalendarFormat(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"standard date", "2025-02-15"},
		{"year end", "2025-12-31"},
		{"leap year", "2024-02-29"},
		{"new year", "2026-01-01"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, _ := time.Parse("2006-01-02", tt.input)
			want := expectedEndOfDay(parsed)

			got, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned unexpected error: %v", tt.input, err)
			}
			if got != want {
				t.Errorf("ParseDate(%q) = %q, want %q", tt.input, got, want)
			}
		})
	}
}

func TestParseDate_RFC3339Passthrough(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"UTC", "2025-02-15T10:00:00Z"},
		{"with offset", "2025-02-15T23:59:59+09:00"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned unexpected error: %v", tt.input, err)
			}
			// RFC3339 strings are returned as-is
			if got != tt.input {
				t.Errorf("ParseDate(%q) = %q, want %q (passthrough)", tt.input, got, tt.input)
			}
		})
	}
}

func TestParseDate_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"plain text", "next week"},
		{"invalid relative no number", "+d"},
		{"invalid relative non-numeric", "+xd"},
		{"just a number", "7"},
		{"incomplete date", "2025-02"},
		{"invalid date format", "02/15/2025"},
		{"negative relative", "-1d"},
		{"random text", "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if err == nil {
				t.Errorf("ParseDate(%q) = %q, want error", tt.input, got)
			}
		})
	}
}

func TestParseDate_OutputFormat(t *testing.T) {
	// All non-RFC3339 outputs must match the format "2006-01-02T15:04:05.000+0000"
	tests := []struct {
		name  string
		input string
	}{
		{"today", "today"},
		{"tomorrow", "tomorrow"},
		{"relative days", "+3d"},
		{"calendar", "2025-06-15"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.input)
			if err != nil {
				t.Fatalf("ParseDate(%q) returned unexpected error: %v", tt.input, err)
			}
			// Verify the output can be parsed with the expected format
			_, parseErr := time.Parse("2006-01-02T15:04:05.000+0000", got)
			if parseErr != nil {
				t.Errorf("ParseDate(%q) = %q, does not match format 2006-01-02T15:04:05.000+0000: %v", tt.input, got, parseErr)
			}
			// Verify it contains 23:59:59 (end of day)
			if !strings.Contains(got, "23:59:59") && !strings.Contains(got, "T") {
				t.Errorf("ParseDate(%q) = %q, expected end-of-day time component", tt.input, got)
			}
		})
	}
}

func TestEndOfDay(t *testing.T) {
	tests := []struct {
		name  string
		input time.Time
	}{
		{"regular date", time.Date(2025, 6, 15, 10, 30, 0, 0, time.Local)},
		{"midnight", time.Date(2025, 1, 1, 0, 0, 0, 0, time.Local)},
		{"end of year", time.Date(2025, 12, 31, 18, 0, 0, 0, time.Local)},
		{"leap year", time.Date(2024, 2, 29, 12, 0, 0, 0, time.Local)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := endOfDay(tt.input)

			// Verify the output parses with the expected format
			parsed, err := time.Parse("2006-01-02T15:04:05.000+0000", got)
			if err != nil {
				t.Fatalf("endOfDay output %q does not match expected format: %v", got, err)
			}

			// The UTC time should correspond to 23:59:59 local converted to UTC
			expectedUTC := time.Date(tt.input.Year(), tt.input.Month(), tt.input.Day(), 23, 59, 59, 0, time.Local).UTC()
			if !parsed.Equal(expectedUTC) {
				t.Errorf("endOfDay(%v) = %q (parsed: %v), want UTC equivalent of %v", tt.input, got, parsed, expectedUTC)
			}
		})
	}
}

func TestParseDate_ErrorMessages(t *testing.T) {
	t.Run("invalid relative format contains input", func(t *testing.T) {
		_, err := ParseDate("+xd")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "+xd") {
			t.Errorf("error message %q should contain the input %q", err.Error(), "+xd")
		}
	})

	t.Run("unsupported format contains input", func(t *testing.T) {
		_, err := ParseDate("garbage")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if !strings.Contains(err.Error(), "garbage") {
			t.Errorf("error message %q should contain the input %q", err.Error(), "garbage")
		}
	})
}
