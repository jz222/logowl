package store

import (
	"context"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type interfaceAnalytics interface {
	InsertOne(models.Analytics) (primitive.ObjectID, error)
	DeleteOne(bson.M) (int64, error)
	DeleteMany(bson.M) (int64, error)
}

type analytics struct {
	db *mongo.Database
}

func (a *analytics) InsertOne(analyticsDocument models.Analytics) (primitive.ObjectID, error) {
	collection := a.db.Collection(CollectionAnalytics)
	res, err := collection.InsertOne(context.TODO(), analyticsDocument)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return res.InsertedID.(primitive.ObjectID), nil
}

func (a *analytics) DeleteOne(filter bson.M) (int64, error) {
	collection := a.db.Collection(CollectionAnalytics)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (a *analytics) DeleteMany(filter bson.M) (int64, error) {
	collection := a.db.Collection(CollectionAnalytics)

	res, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, nil
	}

	return res.DeletedCount, nil
}
