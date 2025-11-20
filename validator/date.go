package validator

import (
	"fmt"
	"time"
)

// parseDateString normalizes an optional date string into a time.Time pointer.
func parseDateString(dateStr *string) (*time.Time, error) {
	if dateStr == nil || *dateStr == "" {
		return nil, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02",
	}

	var lastErr error
	for _, layout := range layouts {
		t, err := time.Parse(layout, *dateStr)
		if err == nil {
			return &t, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("invalid date format: %w", lastErr)
}

// ParseDate converts a date string into a time.Time pointer, supporting multiple layouts.
func ParseDate(dateStr *string) (*time.Time, error) {
	return parseDateString(dateStr)
}
