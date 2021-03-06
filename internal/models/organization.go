package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SlackWebhooks contains information about a Slack webhook.
type SlackWebhooks struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

// Organization contains all properties of an organization.
type Organization struct {
	ID               primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string             `json:"name" bson:"name"`
	Identifier       string             `json:"identifier" bson:"identifier"`
	ReceivedRequests map[string]struct {
		Errors    int `json:"errors" bson:"errors"`
		Analytics int `json:"analytics" bson:"analytics"`
	} `json:"receivedRequests,omitempty" bson:"receivedRequests,omitempty"`
	MonthlyRequestLimit int       `json:"monthlyRequestLimit,omitempty" bson:"monthlyRequestLimit,omitempty"`
	Plan                string    `json:"plan" bson:"plan"`
	SubscriptionID      string    `json:"subscriptionId,omitempty" bson:"subscriptionId,omitempty"`
	PaidThroughDate     string    `json:"paidThroughDate" bson:"paidThroughDate"`
	IsSetUp             bool      `json:"isSetUp" bson:"isSetUp"`
	CreatedAt           time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt           time.Time `json:"upadtedAt" bson:"updatedAt"`
}

// CanBeDeleted determines if the organization can be deleted.
func (o Organization) CanBeDeleted() bool {
	return o.Plan == "free" || o.PaidThroughDate != ""
}

// Validate validates the data of an organization.
func (o Organization) Validate() bool {
	if o.Name == "" {
		return false
	}

	return true
}
