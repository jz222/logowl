package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Adapter struct {
	Name    string `json:"name" bson:"name"`
	Type    string `json:"type" bson:"type"`
	Version string `json:"version" bson:"version"`
}

type Logs struct {
	Timestamp int64  `json:"timestamp", bson:"timestamp"`
	Type      string `json:"type", bson:"type"`
	Log       string `json:"log", bson:"log"`
}

type Error struct {
	ID          *primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Message     string              `json:"message" bson:"message"`
	Stacktrace  string              `json:"stacktrace" bson:"stacktrace"`
	Evolution   map[string]int      `json:"evolution" bson:"evolution"`
	Path        string              `json:"path" bson:"path"`
	Line        string              `json:"line" bson:"line"`
	Type        string              `json:"type" bson:"type"`
	Adapter     Adapter             `json:"adapter" bson:"adapter"`
	Fingerprint string              `json:"fingerprint" bson:"fingerprint"`
	Badges      map[string]string   `json:"badges,omitempty" bson:"badges,omitempty"`
	Snippet     map[string]string   `json:"snippet" bson:"snippet"`
	Logs        []Logs              `json:"logs" bson:"logs"`
	Ticket      string              `json:"ticket" bson:"ticket"`
	Host        string              `json:"host" bson:"host"`
	UserAgent   string              `json:"userAgent" bson:"userAgent"`
	ClientIP    string              `json:"clientIp" bson:"clientIp"`
	Count       int                 `json:"count,omitempty" bson:"count,omitempty"`
	Timestamp   int64               `json:"timestamp" bson:"timestamp"`
	Resolved    bool                `json:"resolved" bson:"resolved"`
	CreatedAt   time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt" bson:"updatedAt"`
}
