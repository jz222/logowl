package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"github.com/mssola/user_agent"
	"go.mongodb.org/mongo-driver/bson"
)

// InterfaceLogging represents the interface
// of the logging service.
type InterfaceLogging interface {
	SaveError(models.Error)
	SaveAnalyticEvent(models.AnalyticEvent)
}

// Logging represents a logging service instance.
type Logging struct {
	Store   store.InterfaceStore
	Request InterfaceRequest
}

// SaveError prepares the error data and saves
// it to the database. If a similar error exits
// already, it updates the existing data to
// achieve aggregation.
func (l *Logging) SaveError(errorEvent models.Error) {
	service, err := l.Store.Service().FindOne(bson.M{"ticket": errorEvent.Ticket})
	if err != nil {
		return
	}

	if errorEvent.Adapter.Type == "browser" {
		ua := user_agent.New(errorEvent.UserAgent)

		osInfo := ua.OS()
		isMobile := ua.Mobile()
		browser, version := ua.Browser()

		errorEvent.Metrics.Platform = osInfo
		errorEvent.Metrics.Browser = fmt.Sprintf("%s %s", browser, version)
		errorEvent.Metrics.IsMobile = strconv.FormatBool(isMobile)
	}

	hash := md5.Sum([]byte(errorEvent.Message + errorEvent.Stacktrace + errorEvent.Ticket))

	dateTool := utils.DateTool{
		Timestamp: errorEvent.Timestamp,
	}

	convertedTimestamp, err := dateTool.GetTimestampBeginnOfDayString()
	if err != nil {
		log.Println("failed to convert timestamp:", errorEvent.Timestamp)
	}

	timestamp := time.Now()

	errorEvent.Fingerprint = hex.EncodeToString(hash[:])
	errorEvent.Evolution = map[string]int{convertedTimestamp: 1}
	errorEvent.Count = 1
	errorEvent.LastSeen = errorEvent.Timestamp
	errorEvent.CreatedAt = timestamp
	errorEvent.UpdatedAt = timestamp

	errorID, err := l.Store.Error().InsertOne(errorEvent)
	if err == nil && service.SlackWebhookURL != "" {
		errorEvent.ID = &errorID
		l.Request.SendSlackAlert(service, errorEvent)
	}
	if err == nil {
		return
	}

	key := fmt.Sprintf("%s.%s", "evolution", convertedTimestamp)

	updatedErrorEvent, err := l.Store.Error().FindOneAndUpdate(
		bson.M{"fingerprint": errorEvent.Fingerprint},
		bson.M{
			"$inc": bson.M{"count": 1, key: 1},
			"$set": bson.M{"lastSeen": errorEvent.Timestamp, "updatedAt": timestamp},
		},
		true,
	)
	if err != nil {
		log.Println(err.Error())
	}

	if service.SlackWebhookURL != "" {
		l.Request.SendSlackAlert(service, updatedErrorEvent)
	}
}

// SaveAnalyticEvent prepares analytic data and saves
// it to the database. Every event is stored in a
// document that represents the statistics of a
// service for the current month.
func (l *Logging) SaveAnalyticEvent(analyticEvent models.AnalyticEvent) {
	_, err := l.Store.Service().CheckPresence(bson.M{"ticket": analyticEvent.Ticket})
	if err != nil {
		return
	}

	const aggregatedMonthlyDataPath = "aggregatedMonthlyData."

	timestamp := time.Now()

	dateTool := utils.DateTool{
		Timestamp: timestamp.Unix(),
	}

	// Get timestamps for the beginn of the hour,
	// beginn of the day and beginn of the month.
	formattedHour, _ := dateTool.GetTimestampBeginnOfHour()
	formattedHourString, _ := dateTool.GetTimestampBeginnOfHourString()
	formattedDay, _ := dateTool.GetTimestampBeginnOfDay()
	formattedMonth, _ := dateTool.GetTimestampBeginnOfMonth()
	humanReadableMonth, _ := dateTool.GetTimestampBeginnOfMonthHumanReadable()

	// Create a prefix for the data that will be written in the document
	// that represents the statistics of the current month.
	prefix := fmt.Sprintf("%s.%s.", "data", formattedHourString)

	// Prepare user agent information
	ua := user_agent.New(analyticEvent.UserAgent)

	isMobile := ua.Mobile()
	browser, _ := ua.Browser()

	incrementUpdate := bson.M{}

	incrementUpdate[prefix+"v"] = 1
	incrementUpdate[aggregatedMonthlyDataPath+"v"] = 1

	incrementUpdate[prefix+"tT"] = analyticEvent.TimeOnPage
	incrementUpdate[aggregatedMonthlyDataPath+"tT"] = analyticEvent.TimeOnPage

	switch browser {
	case "Chrome":
		incrementUpdate[prefix+"c"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"c"] = 1
	case "Safari":
		incrementUpdate[prefix+"s"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"s"] = 1
	case "Opera":
		incrementUpdate[prefix+"o"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"o"] = 1
	case "Firefox":
		incrementUpdate[prefix+"f"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"f"] = 1
	case "Edge":
		incrementUpdate[prefix+"e"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"e"] = 1
	case "IE":
		incrementUpdate[prefix+"i"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"i"] = 1
	default:
		incrementUpdate[prefix+"oB"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"oB"] = 1
	}

	if isMobile {
		incrementUpdate[prefix+"mbl"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"mbl"] = 1
	} else {
		incrementUpdate[prefix+"b"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"b"] = 1
	}

	if analyticEvent.IsNewVisitor {
		incrementUpdate[prefix+"n"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"n"] = 1
	}

	if analyticEvent.IsNewSession {
		incrementUpdate[prefix+"tS"] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"tS"] = 1
	}

	if analyticEvent.Referrer != "" {
		escaped := strings.Replace(analyticEvent.Referrer, ".", "%2E", -1)
		incrementUpdate[prefix+"r."+escaped] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"r."+escaped] = 1
	}

	if analyticEvent.Page != "" {
		escaped := strings.Replace(analyticEvent.Page, ".", "%2E", -1)
		incrementUpdate[prefix+"p."+escaped] = 1
		incrementUpdate[aggregatedMonthlyDataPath+"p."+escaped] = 1
	}

	// Increment existing data or create data in the respective document
	_, err = l.Store.Analytics().FindOneAndUpdate(
		bson.M{"ticket": analyticEvent.Ticket, "month": formattedMonth, "humanReadableMonth": humanReadableMonth},
		bson.M{
			"$inc": incrementUpdate,
			"$set": bson.M{
				fmt.Sprintf("%sday", prefix):  formattedDay,
				fmt.Sprintf("%shour", prefix): formattedHour,
				"updatedAt":                   timestamp,
			},
		},
	)
	if err != nil {
		log.Println(err.Error())
	}
}

// GetLoggingService returns a logging service instance.
func GetLoggingService(store store.InterfaceStore) Logging {
	return Logging{store, &Request{}}
}
