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

func GetLoggingService(store store.InterfaceStore) Logging {
	return Logging{store, &Request{}}
}
