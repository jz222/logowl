package models

type Badge struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Snippet struct {
	Number  int    `json:"number" bson:"number"`
	Snippet string `json:"snippet" bson:"snippet"`
}

type Error struct {
	Message       string    `json:"message" bson:"message"`
	Stacktrace    string    `json:"stacktrace" bson:"stacktrace"`
	Path          string    `json:"path" bson:"path"`
	Line          string    `json:"line" bson:"line"`
	Type          string    `json:"type" bson:"type"`
	Fingerprint   string    `json:"fingerprint" bson:"fingerprint"`
	Badges        []Badge   `json:"badges" bson:"badges"`
	Snippet       []Snippet `json:"snippet" bson:"snippet"`
	Logs          []string  `json:"logs" bson:"logs"`
	Ticket        string    `json:"ticket" bson:"ticket"`
	Host          string    `json:"host" bson:"host"`
	RemoteAddress string    `json:"remoteAddress" bson:"remoteAddress"`
	RequestURL    string    `json:"requestURL" bson:"requestURL"`
	Timestamp     string    `json:"timestamp" bson:"timestamp"`
}
