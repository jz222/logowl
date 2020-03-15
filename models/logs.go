package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Badge struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Error struct {
	Message       string             `json:"message" bson:"message"`
	Stacktrace    string             `json:"stacktrace" bson:"stacktrace"`
	Path          string             `json:"path" bson:"path"`
	Line          string             `json:"line" bson:"line"`
	Type          string             `json:"type" bson:"type"`
	Fingerprint   string             `json:"fingerprint" bson:"fingerprint"`
	Badges        []Badge            `json:"badges" bson:"badges"`
	Snippet       []string           `json:"snippet" bson:"snippet"`
	Logs          []string           `json:"logs" bson:"logs"`
	ProjectID     primitive.ObjectID `json:"projectId,omitempty" bson:"projectId"`
	Host          string             `json:"host" bson:"host"`
	RemoteAddress string             `json:"remoteAddress" bson:"remoteAddress"`
	RequestURL    string             `json:"requestURL" bson:"requestURL"`
	Timestamp     string             `json:"timestamp" bson:"timestamp"`
}
