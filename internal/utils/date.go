package utils

import (
	"strconv"
	"time"
)

// FormatTimestamp returns a new timestamp for the respective day.
func FormatTimestamp(timestamp int64) (int64, string, error) {
	parsed := time.Unix(timestamp, 0)

	day := parsed.Format("2006-01-02")

	formatted, err := time.Parse("2006-01-02", day)
	if err != nil {
		return 0, "", err
	}

	return formatted.Unix(), strconv.FormatInt(formatted.Unix(), 10), nil
}
