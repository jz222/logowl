package services

import (
	"context"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type event struct {
	DB *mongo.Database
}

func (e *event) DeleteError(filter bson.M) (int64, error) {
	collection := mongodb.GetClient().Collection(mongodb.Errors)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (e *event) GetError(filter bson.M) (*models.Error, error) {
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

func (e *event) GetErrors(ticket string, page int64) (*[]models.Error, error) {
	collection := mongodb.GetClient().Collection(mongodb.Errors)

	cur, err := collection.Find(
		context.TODO(),
		bson.M{"ticket": ticket},
		options.MergeFindOptions().SetSort(bson.M{"updatedAt": -1}),
		options.MergeFindOptions().SetSkip(page*10),
		options.MergeFindOptions().SetLimit(10),
	)
	if err != nil {
		return nil, err
	}

	var errorEvents []models.Error

	for cur.Next(context.TODO()) {
		var errorEvent models.Error

		err := cur.Decode(&errorEvent)
		if err == nil {
			errorEvents = append(errorEvents, errorEvent)
		}
	}

	if errorEvents == nil {
		errorEvents = []models.Error{}
	}

	return &errorEvents, nil
}

func (e *event) UpdateError(filter, update bson.M) error {
	collection := mongodb.GetClient().Collection(mongodb.Errors)

	update["updatedAt"] = time.Now()

	res := collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func GetEventService(db *mongo.Database) event {
	return event{
		DB: db,
	}
}
