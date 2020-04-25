package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Adapter contains information about the adapter.
type Adapter struct {
	Name    string `json:"name" bson:"name"`
	Type    string `json:"type" bson:"type"`
	Version string `json:"version" bson:"version"`
}

// Logs contains all properties of a log.
type Logs struct {
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Type      string `json:"type" bson:"type"`
	Log       string `json:"log" bson:"log"`
}

// Metrics contains information about the system
type Metrics struct {
	Platform string `json:"platform,omitempty" bson:"platform,omitempty"`
	Browser  string `json:"browser,omitempty" bson:"browser,omitempty"`
	IsMobile string `json:"isMobile,omitempty" bson:"isMobile,omitempty"`
}

// UserInteraction contains information about an element that was clicked by the user.
type UserInteraction struct {
	Timestamp int64  `json:"timestamp" bson:"timestamp"`
	Element   string `json:"element" bson:"element"`
	InnerText string `json:"innerText" bson:"innerText"`
	ElementID string `json:"elementId" bson:"elementId"`
	Location  string `json:"location" bson:"location"`
}

// Error contains all properties of an error event.
type Error struct {
	ID               *primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Message          string               `json:"message" bson:"message"`
	Stacktrace       string               `json:"stacktrace" bson:"stacktrace"`
	Evolution        map[string]int       `json:"evolution" bson:"evolution"`
	Path             string               `json:"path" bson:"path"`
	Line             string               `json:"line" bson:"line"`
	Type             string               `json:"type" bson:"type"`
	Adapter          Adapter              `json:"adapter" bson:"adapter"`
	Fingerprint      string               `json:"fingerprint" bson:"fingerprint"`
	Badges           map[string]string    `json:"badges,omitempty" bson:"badges,omitempty"`
	Snippet          map[string]string    `json:"snippet,omitempty" bson:"snippet,omitempty"`
	Logs             []Logs               `json:"logs,omitempty" bson:"logs,omitempty"`
	Ticket           string               `json:"ticket" bson:"ticket"`
	Host             string               `json:"host,omitempty" bson:"host,omitempty"`
	UserAgent        string               `json:"userAgent" bson:"userAgent"`
	Metrics          Metrics              `json:"metrics" bson:"metrics"`
	UserInteractions []UserInteraction    `json:"userInteractions,omitempty" bson:"userInteractions,omitempty"`
	ClientIP         string               `json:"clientIp" bson:"clientIp"`
	Count            int                  `json:"count,omitempty" bson:"count,omitempty"`
	Timestamp        int64                `json:"timestamp" bson:"timestamp"`
	Resolved         bool                 `json:"resolved" bson:"resolved"`
	SeenBy           []primitive.ObjectID `json:"seenBy,omitempty" bson:"seenBy,omitempty"`
	LastSeen         int64                `json:"lastSeen" bson:"lastSeen"`
	CreatedAt        time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time            `json:"updatedAt" bson:"updatedAt"`
}

func (e *Error) IsValid() bool {
	if len(e.Logs) > 50 {
		return false
	}

	if len(e.UserInteractions) > 50 {
		return false
	}

	if len(e.SeenBy) > 0 {
		return false
	}

	if len(e.Badges) > 200 {
		return false
	}

	if len(e.Snippet) > 50 {
		return false
	}

	if len(e.Evolution) > 0 {
		return false
	}

	return true
}

type AnalyticData struct {
	Day             int64          `json:"day" bson:"day"`
	Windows         int            `json:"windows" bson:"wndws,omitempty"`
	Mac             int            `json:"mc" bson:"mc,omitempty"`
	Linux           int            `json:"linux" bson:"lnx,omitempty"`
	OtherPlatforms  int            `json:"otherPlatforms" bson:"othrPltfrms,omitempty"`
	Chrome          int            `json:"chrome" bson:"chrm,omitempty"`
	Firefox         int            `json:"firefox" bson:"frfx,omitempty"`
	Safari          int            `json:"safari" bson:"sfr,omitempty"`
	Edge            int            `json:"edge" bson:"edg,omitempty"`
	IE              int            `json:"ie" bson:"ie,omitempty"`
	Opera           int            `json:"opera" bson:"opr,omitempty"`
	OtherBrowsers   int            `json:"otherBrowsers" bson:"othrBrwsrs,omitempty"`
	Mobile          int            `json:"mobile" bson:"mbl,omitempty"`
	Tablet          int            `json:"tablet" bson:"tblt,omitempty"`
	Browser         int            `json:"browser" bson:"brwsr,omitempty"`
	Visits          int            `json:"visits" bson:"vsts,omitempty"`
	NewVisitors     int            `json:"newVisitors" bson:"nwVstrs,omitempty"`
	TotalSessions   int            `json:"totalSessions" bson:"ttlSssns,omitempty"`
	TotalTimeOnPage int            `json:"totalTimeOnPage" bson:"ttlTmOnPg,omitempty"`
	EntryPage       map[string]int `json:"entryPage" bson:"entryPg,omitempty"`
	Referrer        map[string]int `json:"referrer" bson:"rfrr,omitempty"`
}

type Analytics struct {
	Ticket             string                  `json:"ticket" bson:"ticket"`
	Month              int64                   `json:"month" bson:"month"`
	HumanReadableMonth string                  `json:"humanReadableMonth" bson:"humanReadableMonth"`
	Data               map[string]AnalyticData `json:"data" bson:"data"`
	CreatedAt          time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt          time.Time               `json:"updatedAt" bson:"updatedAt"`
}

// AnalyticEvent contains information about a page visitor.
type AnalyticEvent struct {
	Ticket       string `json:"ticket" bson:"ticket"`
	IsNewVisitor bool   `json:"isNewVisitor" bson:"isNewVisitor"`
	IsNewSession bool   `json:"isNewSession" bson:"isNewSession"`
	TimeOnPage   int    `json:"timeOnPage" bson:"timeOnPage"`
	Referrer     string `json:"referrer" bson:"referrer"`
	EntryPage    string `json:"entryPage" bson:"entryPage"`
	UserAgent    string `json:"userAgent" bson:"userAgent"`
}

type AnalyticsInsightsPageViews struct {
	Day         int64  `json:"-"`
	Unit        string `json:"unit"`
	Sessions    int    `json:"sessions"`
	Visits      int    `json:"visits"`
	NewVisitors int    `json:"newVisitors"`
}

type AnalyticInsights struct {
	TotalVisits      int                          `json:"totalVisits"`
	TotalNewVisitors int                          `json:"totalNewVisitors"`
	TotalSessions    int                          `json:"totalSessions"`
	Data             []AnalyticsInsightsPageViews `json:"pageViews"`
}
