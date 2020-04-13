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
	Timestamp int64  `json:"timestamp", bson:"timestamp"`
	Type      string `json:"type", bson:"type"`
	Log       string `json:"log", bson:"log"`
}

// Metrics contains information about the system
type Metrics struct {
	Platform string `json:"platform,omitempty" bson:"platform,omitempty"`
	Browser  string `json:"browser,omitempty" bson:"browser,omitempty"`
	IsMobile string `json:"isMobile,omitempty" bson:"isMobile,omitempty"`
}

// Error contains all properties of an error event.
type Error struct {
	ID          *primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	Message     string               `json:"message" bson:"message"`
	Stacktrace  string               `json:"stacktrace" bson:"stacktrace"`
	Evolution   map[string]int       `json:"evolution" bson:"evolution"`
	Path        string               `json:"path" bson:"path"`
	Line        string               `json:"line" bson:"line"`
	Type        string               `json:"type" bson:"type"`
	Adapter     Adapter              `json:"adapter" bson:"adapter"`
	Fingerprint string               `json:"fingerprint" bson:"fingerprint"`
	Badges      map[string]string    `json:"badges,omitempty" bson:"badges,omitempty"`
	Snippet     map[string]string    `json:"snippet,omitempty" bson:"snippet,omitempty"`
	Logs        []Logs               `json:"logs,omitempty" bson:"logs,omitempty"`
	Ticket      string               `json:"ticket" bson:"ticket"`
	Host        string               `json:"host" bson:"host"`
	UserAgent   string               `json:"userAgent" bson:"userAgent"`
	Metrics     Metrics              `json:"metrics" bson:"metrics"`
	ClientIP    string               `json:"clientIp" bson:"clientIp"`
	Count       int                  `json:"count,omitempty" bson:"count,omitempty"`
	Timestamp   int64                `json:"timestamp" bson:"timestamp"`
	Resolved    bool                 `json:"resolved" bson:"resolved"`
	SeenBy      []primitive.ObjectID `json:"seenBy,omitempty" bson:"seenBy,omitempty"`
	LastSeen    int64                `json:"lastSeen" bson:"lastSeen"`
	CreatedAt   time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt" bson:"updatedAt"`
}
