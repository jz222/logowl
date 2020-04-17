package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Service contains all the properties of a service.
type Service struct {
	ID              primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name            string             `json:"name" bson:"name"`
	Description     string             `json:"description" bson:"description"`
	Type            string             `json:"type" bson:"type"`
	OrganizationID  primitive.ObjectID `json:"organizationId" bson:"organizationId"`
	Ticket          string             `json:"ticket" bson:"ticket"`
	SlackWebhookURL string             `json:"slackWebhookURL,omitempty" bson:"slackWebhookURL,omitempty"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
}

// Validate validates the data of a service.
func (s *Service) Validate() bool {
	if s.Name == "" || s.Type == "" || s.Description == "" || s.OrganizationID.IsZero() {
		return false
	}

	return true
}
