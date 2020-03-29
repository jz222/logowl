package event

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func GetError(filter bson.M) (*models.Error, error) {
	var errorEvent models.Error

	collection := mongodb.GetClient().Collection(mongodb.Errors)

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return nil, queryResult.Err()
	}

	err := queryResult.Decode(&errorEvent)
	if err != nil {
		return nil, err
	}

	return &errorEvent, nil
}
