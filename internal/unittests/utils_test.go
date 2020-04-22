package unittests

import (
	"testing"

	"github.com/jz222/loggy/internal/utils"
)

func TestFormatTimestamp(t *testing.T) {
	convertedTimestamp, convertedTimestampString, err := utils.FormatTimestampToBeginnOfDay(1585930192)
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

func TestGenerateTicket(t *testing.T) {
	ticket, err := utils.GenerateTicket()
	if err != nil {
		t.Error(err.Error())
	}

	if len(ticket) != 50 {
		t.Errorf("Expected ticket to have %d characters, got %d", 50, len(ticket))
	}
}
