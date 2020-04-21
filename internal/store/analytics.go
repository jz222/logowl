package store

import (
	"context"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type interfaceAnalytics interface {
	InsertOne(models.Analytics) (primitive.ObjectID, error)
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
