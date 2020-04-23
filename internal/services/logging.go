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

	_, convertedTimestamp, err := utils.FormatTimestampToBeginnOfDay(errorEvent.Timestamp)
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

	timestamp := time.Now()

	_, formattedTs, err := utils.FormatTimestampToHour(timestamp.Unix())
	if err != nil {
		return
	}

	formattedMonth, _, humanReadableMonth, err := utils.FormatTimestampToMonth(timestamp.Unix())
	if err != nil {
		return
	}

	prefix := fmt.Sprintf("%s.%s.", "data", formattedTs)

	ua := user_agent.New(analyticEvent.UserAgent)

	isMobile := ua.Mobile()
	browser, _ := ua.Browser()

	incrementUpdate := bson.M{}

	incrementUpdate[prefix+"vsts"] = 1
	incrementUpdate[prefix+"ttlTmOnPg"] = analyticEvent.TimeOnPage

	switch browser {
	case "Chrome":
		incrementUpdate[prefix+"chrm"] = 1
	case "Safari":
		incrementUpdate[prefix+"sfr"] = 1
	case "Opera":
		incrementUpdate[prefix+"opr"] = 1
	case "Firefox":
		incrementUpdate[prefix+"frfx"] = 1
	case "Edge":
		incrementUpdate[prefix+"edg"] = 1
	case "IE":
		incrementUpdate[prefix+"ie"] = 1
	default:
		incrementUpdate[prefix+"othrBrwsrs"] = 1
	}

	if isMobile {
		incrementUpdate[prefix+"mbl"] = 1
	} else {
		incrementUpdate[prefix+"brwsr"] = 1
	}

	if analyticEvent.IsNewVisitor {
		incrementUpdate[prefix+"nwVstrs"] = 1
	}

	if analyticEvent.IsNewSession {
		incrementUpdate[prefix+"ttlSssns"] = 1
	}

	if analyticEvent.Referrer != "" {
		escaped := strings.Replace(analyticEvent.Referrer, ".", "%2E", -1)
		incrementUpdate[prefix+"rfrr."+escaped] = 1
	}

	if analyticEvent.EntryPage != "" {
		escaped := strings.Replace(analyticEvent.EntryPage, ".", "%2E", -1)
		incrementUpdate[prefix+"entryPg."+escaped] = 1
	}

	_, err = l.Store.Analytics().FindOneAndUpdate(
		bson.M{"ticket": analyticEvent.Ticket, "month": formattedMonth, "humanReadableMonth": humanReadableMonth},
		bson.M{
			"$inc": incrementUpdate,
			"$set": bson.M{"updatedAt": timestamp},
		},
	)
	if err != nil {
		log.Println(err.Error())
	}
}

func GetLoggingService(store store.InterfaceStore) Logging {
	return Logging{store, &Request{}}
}
