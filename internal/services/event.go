package services

import (
	"errors"
	"sort"
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
	DeleteErrors(bson.M) (int64, error)
	UpdateError(bson.M, bson.M) error
	GetAnalytics(string, string) (models.AnalyticInsights, error)
}

type Event struct {
	Store store.InterfaceStore
}

func (e *Event) DeleteError(filter bson.M) (int64, error) {
	return e.Store.Error().DeleteOne(filter)
}

func (e *Event) DeleteErrors(filter bson.M) (int64, error) {
	return e.Store.Error().DeleteMany(filter)
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
	if mode != "today" && mode != "lastSevenDays" && mode != "lastFourteenDays" && mode != "lastMonth" {
		return models.AnalyticInsights{}, errors.New("the provided mode is invalid")
	}

	var timeframeStart int64
	var timeframeEnd int64

	var response models.AnalyticInsights

	dateTool := utils.DateTool{
		Timestamp: time.Now().Unix(),
	}

	filter := bson.M{"ticket": ticket}

	// Calculate the timeframe for the current day
	if mode == "today" {
		currentMonth, _ := dateTool.GetTimestampBeginnOfMonth()
		filter["month"] = currentMonth

		startTime, endTime, _ := dateTool.GetTimeframeToday()
		timeframeStart = startTime
		timeframeEnd = endTime
	}

	// Calculate the timeframe for the last seven days
	if mode == "lastSevenDays" {
		previousMonth, _ := dateTool.GetTimestampBeginnOfPreviousMonth()
		filter["month"] = bson.M{"$gte": previousMonth}

		startTime, endTime, _ := dateTool.GetTimeframeLastSevenDays()
		timeframeStart = startTime
		timeframeEnd = endTime
	}

	// Calculate the timeframe for the last fourteen days
	if mode == "lastFourteenDays" {
		previousMonth, _ := dateTool.GetTimestampBeginnOfPreviousMonth()
		filter["month"] = bson.M{"$gte": previousMonth}

		startTime, endTime, _ := dateTool.GetTimeframeLastFourteenDays()
		timeframeStart = startTime
		timeframeEnd = endTime
	}

	// Calculate the timeframe for the last month
	if mode == "lastMonth" {
		previousMonth, _ := dateTool.GetTimestampBeginnOfPreviousMonth()
		filter["month"] = bson.M{"$gte": previousMonth}

		startTime, endTime, _ := dateTool.GetTimeframeLastMonth()
		timeframeStart = startTime
		timeframeEnd = endTime
	}

	// Find all documents in the database that match the filter
	analyticDocuments, err := e.Store.Analytics().Find(filter)
	if err != nil {
		return models.AnalyticInsights{}, err
	}

	// Filter all the days that are within the timeframe
	for _, analyticDocument := range analyticDocuments {
		for _, v := range analyticDocument.Data {
			if v.Day >= timeframeStart && v.Day <= timeframeEnd {
				response.Data = append(response.Data, v)
			}
		}
	}

	// Sort the results either by hour if the mode
	// is set to today or by day for any other mode.
	if mode == "today" {
		sort.Slice(response.Data, func(i, j int) bool {
			return response.Data[i].Hour < response.Data[j].Hour
		})
	} else {
		sort.Slice(response.Data, func(i, j int) bool {
			return response.Data[i].Day < response.Data[j].Day
		})
	}

	var currentDay int64
	var aggregatedData []models.AnalyticData

	totalVisits := 0
	totalNewVisitors := 0
	totalSessions := 0

	// Aggregate the data either by hour or by day
	for _, metrics := range response.Data {
		totalVisits += metrics.Visits
		totalNewVisitors += metrics.NewVisitors
		totalSessions += metrics.TotalSessions

		if mode == "today" {
			continue
		}

		if currentDay != metrics.Day {
			currentDay = metrics.Day
			aggregatedData = append(aggregatedData, metrics)
			continue
		}

		prevIndex := len(aggregatedData) - 1

		// Aggregate flat data
		aggregatedData[prevIndex].NewVisitors += metrics.NewVisitors
		aggregatedData[prevIndex].TotalSessions += metrics.TotalSessions
		aggregatedData[prevIndex].Visits += metrics.Visits
		aggregatedData[prevIndex].Chrome += metrics.Chrome
		aggregatedData[prevIndex].Firefox += metrics.Firefox
		aggregatedData[prevIndex].Safari += metrics.Safari
		aggregatedData[prevIndex].Opera += metrics.Opera
		aggregatedData[prevIndex].Edge += metrics.Edge
		aggregatedData[prevIndex].IE += metrics.IE
		aggregatedData[prevIndex].OtherBrowsers += metrics.OtherBrowsers
		aggregatedData[prevIndex].Mobile += metrics.Mobile
		aggregatedData[prevIndex].Tablet += metrics.Tablet
		aggregatedData[prevIndex].Desktop += metrics.Desktop
		aggregatedData[prevIndex].TotalTimeOnPage += metrics.TotalTimeOnPage

		// Aggregate referrers map
		for k, v := range metrics.Referrer {
			if _, ok := aggregatedData[prevIndex].Referrer[k]; ok {
				aggregatedData[prevIndex].Referrer[k] += v
			} else {
				if aggregatedData[prevIndex].Referrer == nil {
					aggregatedData[prevIndex].Referrer = map[string]int{}
				}
				aggregatedData[prevIndex].Referrer[k] = v
			}
		}

		// Aggregate visited pages map
		for k, v := range metrics.Pages {
			if _, ok := aggregatedData[prevIndex].Pages[k]; ok {
				aggregatedData[prevIndex].Pages[k] += v
			} else {
				if aggregatedData[prevIndex].Pages == nil {
					aggregatedData[prevIndex].Pages = map[string]int{}
				}
				aggregatedData[prevIndex].Pages[k] = v
			}
		}
	}

	if mode != "today" {
		response.Data = aggregatedData
	}

	response.TimeframeStart = timeframeStart
	response.TimeframeEnd = timeframeEnd
	response.TotalVisits = totalVisits
	response.TotalNewVisitors = totalNewVisitors
	response.TotalSessions = totalSessions

	return response, nil
}

func GetEventService(store store.InterfaceStore) Event {
	return Event{store}
}
