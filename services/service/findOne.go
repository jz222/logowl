package service

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FindOne(filter bson.M) (models.Service, error) {
	var service models.Service

	collection := mongodb.GetClient().Collection("services")

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return models.Service{}, queryResult.Err()
	}

	err := queryResult.Decode(&service)
	if err != nil {
		return models.Service{}, err
	}

	return service, nil
}
