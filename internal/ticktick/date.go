package ticktick

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseDate parses a human-friendly date string into RFC3339 format.
// Supported formats:
//   - "today", "tomorrow"
//   - "+3d" (relative days)
//   - "YYYY-MM-DD"
//   - Already RFC3339 formatted strings
func ParseDate(s string) (string, error) {
	now := time.Now()

	switch strings.ToLower(s) {
	case "today":
		return endOfDay(now), nil
	case "tomorrow":
		return endOfDay(now.AddDate(0, 0, 1)), nil
	}

	// +Nd format
	if strings.HasPrefix(s, "+") && strings.HasSuffix(s, "d") {
		days, err := strconv.Atoi(s[1 : len(s)-1])
		if err != nil {
			return "", fmt.Errorf("invalid relative date format: %s", s)
		}
		return endOfDay(now.AddDate(0, 0, days)), nil
	}

	// YYYY-MM-DD
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return endOfDay(t), nil
	}

	// Already RFC3339
	if _, err := time.Parse(time.RFC3339, s); err == nil {
		return s, nil
	}

	return "", fmt.Errorf("unsupported date format: %s (use today, tomorrow, +Nd, or YYYY-MM-DD)", s)
}

func endOfDay(t time.Time) string {
	eod := time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.Local)
	return eod.UTC().Format("2006-01-02T15:04:05.000+0000")
}
