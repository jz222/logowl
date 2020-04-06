package store

import (
	"context"

	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interfaceErrorEvent interface {
	DeleteOne(bson.M) (int64, error)
	DeleteMany(bson.M) (int64, error)
	FindOne(bson.M) (models.Error, error)
	FindOneAndUpdate(bson.M, bson.M, bool) error
	FindPaged(bson.M, int64) ([]models.Error, error)
	InsertOne(models.Error) error
}

type errorEvent struct {
	db *mongo.Database
}

func (e *errorEvent) DeleteOne(filter bson.M) (int64, error) {
	collection := e.db.Collection(CollectionErrors)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (e *errorEvent) DeleteMany(filter bson.M) (int64, error) {
	collection := e.db.Collection(CollectionErrors)

	res, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (e *errorEvent) FindOne(filter bson.M) (models.Error, error) {
	var errorEvent models.Error

	collection := e.db.Collection(CollectionErrors)

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return models.Error{}, queryResult.Err()
	}

	err := queryResult.Decode(&errorEvent)
	if err != nil {
		return models.Error{}, err
	}

	return errorEvent, nil
}

func (e *errorEvent) FindPaged(filter bson.M, page int64) ([]models.Error, error) {
	collection := e.db.Collection(CollectionErrors)

	cur, err := collection.Find(
		context.TODO(),
		filter,
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

	return errorEvents, nil
}

func (e *errorEvent) FindOneAndUpdate(filter, update bson.M, upsert bool) error {
	collection := e.db.Collection(CollectionErrors)

	res := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
	)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}

func (e *errorEvent) InsertOne(errorEvent models.Error) error {
	collection := e.db.Collection(CollectionErrors)
	_, err := collection.InsertOne(context.TODO(), errorEvent)
	if err != nil {
		return err
	}

	return nil
}
