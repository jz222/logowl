package services

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"github.com/mssola/user_agent"
	"go.mongodb.org/mongo-driver/bson"
)

type InterfaceLogging interface {
	SaveError(models.Error)
	SaveAnalyticEvent(models.AnalyticEvent)
}

type Logging struct {
	Store   store.InterfaceStore
	Request InterfaceRequest
}

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

	_, convertedTimestamp, err := utils.FormatTimestamp(errorEvent.Timestamp)
	if err != nil {
		log.Println("failed to convert timestamp:", errorEvent.Timestamp)
		return
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

func (l *Logging) SaveAnalyticEvent(analyticEvent models.AnalyticEvent) {
	_, err := l.Store.Service().CheckPresence(bson.M{"ticket": analyticEvent.Ticket})
	if err != nil {
		return
	}

	ua := user_agent.New(analyticEvent.UserAgent)

	//osInfo := ua.OS()
	isMobile := ua.Mobile()
	browser, _ := ua.Browser()

	timestamp := time.Now()

	_, formattedTs, err := utils.FormatTimestampToHour(timestamp.Unix())
	if err != nil {
		return
	}

	analyticData := models.AnalyticData{
		Visitors:        1,
		UniqueVisitors:  1,
		TotalTimeOnPage: analyticEvent.TimeOnPage,
	}

	switch browser {
	case "Chrome":
		analyticData.Chrome = 1
	case "Safari":
		analyticData.Safari = 1
	case "Opera":
		analyticData.Opera = 1
	case "Edge":
		analyticData.Edge = 1
	case "IE":
		analyticData.IE = 1
	default:
		analyticData.OtherBrowsers = 1
	}

	if isMobile {
		analyticData.Mobile = 1
	} else {
		analyticData.Browser = 1
	}

	if analyticEvent.Referrer != "" {
		analyticData.Referrer = map[string]int{analyticEvent.Referrer: 1}
	}

	analyticDocument := models.Analytics{
		Ticket:    analyticEvent.Ticket,
		Data:      map[string]models.AnalyticData{formattedTs: analyticData},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	_, err = l.Store.Analytics().InsertOne(analyticDocument)
	if err != nil {
		fmt.Println(err.Error())
	}
}

func GetLoggingService(store store.InterfaceStore) Logging {
	return Logging{store, &Request{}}
}
