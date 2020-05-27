package unittests

import (
	"testing"

	"github.com/jz222/logowl/internal/utils"
)

func TestFormatTimestamp(t *testing.T) {
	dateTool := utils.DateTool{
		Timestamp: 1585930192,
	}

	convertedTimestamp, err := dateTool.GetTimestampBeginnOfDay()
	if err != nil {
		t.Errorf(err.Error())
	}

	convertedTimestampString, err := dateTool.GetTimestampBeginnOfDayString()
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
