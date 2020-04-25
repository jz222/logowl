package utils

import (
	"strconv"
	"time"
)

type DateTool struct {
	Timestamp int64
}

// GetTimestampBeginnOfDay returns a new timestamp for the respective day.
func (d DateTool) GetTimestampBeginnOfDay() (int64, error) {
	parsed := time.Unix(d.Timestamp, 0)

	day := parsed.Format("2006-01-02")

	formatted, err := time.Parse("2006-01-02", day)
	if err != nil {
		return 0, err
	}

	return formatted.Unix(), nil
}

// GetTimestampBeginnOfDayString returns a new timestamp for the respective day as string.
func (d DateTool) GetTimestampBeginnOfDayString() (string, error) {
	parsed := time.Unix(d.Timestamp, 0)

	day := parsed.Format("2006-01-02")

	formatted, err := time.Parse("2006-01-02", day)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(formatted.Unix(), 10), nil
}

func (d DateTool) GetTimestampBeginnOfMonth() (int64, error) {
	parsed := time.Unix(d.Timestamp, 0)

	month := parsed.Format("2006-01")

	formatted, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, err
	}

	currentMonth := formatted.Unix()

	return currentMonth, nil
}

func (d DateTool) GetTimestampBeginnOfMonthHumanReadable() (string, error) {
	parsed := time.Unix(d.Timestamp, 0)

	month := parsed.Format("2006-01")

	formatted, err := time.Parse("2006-01", month)
	if err != nil {
		return "", err
	}

	currentMonthHumanReadable := formatted.Format("January 2006")

	return currentMonthHumanReadable, nil
}

func (d DateTool) GetTimestampBeginnOfPreviousMonth() (int64, error) {
	parsed := time.Unix(d.Timestamp, 0)

	month := parsed.Format("2006-01")

	formatted, err := time.Parse("2006-01", month)
	if err != nil {
		return 0, err
	}

	previousMonth := formatted.AddDate(0, -1, 0).Unix()

	return previousMonth, nil
}

func (d DateTool) GetTimestampBeginnOfHourString() (string, error) {
	parsed := time.Unix(d.Timestamp, 0)

	hour := parsed.Format("2006-01-02 15")

	formatted, err := time.Parse("2006-01-02 15", hour)
	if err != nil {
		return "", err
	}

	return strconv.FormatInt(formatted.Unix(), 10), nil
}
