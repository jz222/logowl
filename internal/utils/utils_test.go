package utils

import (
	"testing"
)

func TestFormatTimestamp(t *testing.T) {
	convertedTimestamp, convertedTimestampString, err := FormatTimestamp(1585930192)
	if err != nil {
		t.Errorf(err.Error())
	}

	if convertedTimestampString != "1585872000" {
		t.Errorf("Timestamp as string was incorrect, got: %s, expected: %s", convertedTimestampString, "1585872000")
	}

	if convertedTimestamp != 1585872000 {
		t.Errorf("Timestamp was incorred, got: %d, expected: %d", convertedTimestamp, 1585872000)
	}
}
