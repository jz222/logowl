package services

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterfaceEvent interface {
	GetError(bson.M, primitive.ObjectID) (models.Error, error)
	GetErrors(string, int64) ([]models.Error, error)
	DeleteError(bson.M) (int64, error)
	UpdateError(bson.M, bson.M) error
	GetAnalytics(string, string) (models.AnalyticInsights, error)
}

type Event struct {
	Store store.InterfaceStore
}

func (e *Event) DeleteError(filter bson.M) (int64, error) {
	return e.Store.Error().DeleteOne(filter)
}

func (e *Event) GetError(filter bson.M, viewer primitive.ObjectID) (models.Error, error) {
	return e.Store.Error().FindOneAndUpdate(filter, bson.M{"$addToSet": bson.M{"seenBy": viewer}}, true)
}

func (e *Event) GetErrors(ticket string, page int64) ([]models.Error, error) {
	return e.Store.Error().FindPaged(bson.M{"ticket": ticket}, page)
}

func (e *Event) UpdateError(filter, update bson.M) error {
	update["updatedAt"] = time.Now()

	_, err := e.Store.Error().FindOneAndUpdate(filter, bson.M{"$set": update}, false)
	if err != nil {
		return err
	}

	return nil
}

func (e *Event) GetAnalytics(ticket, mode string) (models.AnalyticInsights, error) {
	if mode != "today" && mode != "lastSevenDays" {
		return models.AnalyticInsights{}, errors.New("the provided mode is invalid")
	}

	var timeframeStart int64
	var timeframeEnd int64

	var response models.AnalyticInsights

	timestamp := time.Now()

	filter := bson.M{"ticket": ticket}

	if mode == "today" {
		currentMonth, _, _, _ := utils.FormatTimestampToMonth(timestamp.Unix())
		filter["month"] = currentMonth

		currentDay, _, _ := utils.FormatTimestampToBeginnOfDay(timestamp.Unix())
		timeframeStart = currentDay
		timeframeEnd = currentDay + int64(60*60*24-1)
	}

	if mode == "lastSevenDays" {
		_, _, previousMonth, _ := utils.FormatTimestampToMonth(timestamp.Unix())
		filter["month"] = bson.M{"$gte": previousMonth}

		currentDay, _, _ := utils.FormatTimestampToBeginnOfDay(timestamp.Unix())
		timeframeStart = currentDay - int64(60*60*24*6)
		timeframeEnd = currentDay + int64(60*60*24-1)
	}

	analyticDocuments, err := e.Store.Analytics().Find(filter)
	if err != nil {
		return models.AnalyticInsights{}, err
	}

	for _, analyticDocument := range analyticDocuments {
		for k, v := range analyticDocument.Data {
			parsedKey, err := strconv.ParseInt(k, 10, 64)
			if err != nil {
				continue
			}

			if parsedKey >= timeframeStart && parsedKey <= timeframeEnd {
				metrics := models.AnalyticsInsightsPageViews{
					Day:         v.Day,
					Unit:        k,
					Sessions:    v.TotalSessions,
					Visits:      v.Visits,
					NewVisitors: v.NewVisitors,
				}

				response.Data = append(response.Data, metrics)
			}
		}
	}

	sort.Slice(response.Data, func(i, j int) bool {
		return response.Data[i].Unit < response.Data[j].Unit
	})

	var currentDay int64
	var aggregatedData []models.AnalyticsInsightsPageViews

	totalVisits := 0
	totalNewVisitors := 0
	totalSessions := 0

	for _, metrics := range response.Data {
		totalVisits += metrics.Visits
		totalNewVisitors += metrics.NewVisitors
		totalSessions += metrics.Sessions

		if mode == "today" {
			continue
		}

		if currentDay != metrics.Day {
			currentDay = metrics.Day
			aggregatedData = append(aggregatedData, metrics)
			continue
		}

		prevIndex := len(aggregatedData) - 1

		aggregatedData[prevIndex].NewVisitors += metrics.NewVisitors
		aggregatedData[prevIndex].Sessions += metrics.Sessions
		aggregatedData[prevIndex].Visits += metrics.Visits
	}

	if mode != "today" {
		response.Data = aggregatedData
	}

	response.TotalVisits = totalVisits
	response.TotalNewVisitors = totalNewVisitors
	response.TotalSessions = totalSessions

	return response, nil
}

func GetEventService(store store.InterfaceStore) Event {
	return Event{store}
}
