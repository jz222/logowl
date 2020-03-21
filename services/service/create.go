package service

import (
	"context"
	"errors"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/organization"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(service models.Service) (primitive.ObjectID, error) {
	timestamp := time.Now()
	service.CreatedAt = timestamp
	service.UpdatedAt = timestamp

	if !service.Validate() {
		return primitive.ObjectID{}, errors.New("the provided service data is invalid")
	}

	organizationExists, err := organization.CheckPresence(bson.M{"_id": service.OrganizationID})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if !organizationExists {
		return primitive.ObjectID{}, errors.New("the provided organization does not exist")
	}

	ticket, err := utils.GenerateTicket()
	if err != nil {
		return primitive.ObjectID{}, err
	}

	service.Ticket = ticket

	collection := mongodb.GetClient().Collection("services")

	result, err := collection.InsertOne(context.TODO(), service)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving service to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
