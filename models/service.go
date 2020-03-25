package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Service struct {
	ID             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name" bson:"name"`
	Description    string             `json:"description" bson:"description"`
	Type           string             `json:"type" bson:"type"`
	OrganizationID primitive.ObjectID `json:"organizationId" bson:"organizationId"`
	Ticket         string             `json:"ticket" bson:"ticket"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

func (s *Service) Validate() bool {
	if s.Name == "" || s.Type == "" || s.Description == "" || s.OrganizationID.IsZero() {
		return false
	}

	return true
}
