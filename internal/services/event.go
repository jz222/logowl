package services

import (
	"errors"
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
		timeframeEnd = currentDay + int64(60*60*24)
	}

	if mode == "lastSevenDays" {
		_, _, previousMonth, _ := utils.FormatTimestampToMonth(timestamp.Unix())
		filter["month"] = bson.M{"$gte": previousMonth}

		currentDay, _, _ := utils.FormatTimestampToBeginnOfDay(timestamp.Unix())
		timeframeStart = currentDay - int64(60*60*24*6)
		timeframeEnd = currentDay + int64(60*60*24)
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
				response.PageViews.Sessions = append(response.PageViews.Sessions, v.TotalSessions)
				response.PageViews.Visits = append(response.PageViews.Visits, v.Visits)
				response.PageViews.NewVisitors = append(response.PageViews.NewVisitors, v.NewVisitors)
				response.PageViews.Labels = append(response.PageViews.Labels, k)
			}
		}
	}

	return response, nil
}

func GetEventService(store store.InterfaceStore) Event {
	return Event{store}
}
