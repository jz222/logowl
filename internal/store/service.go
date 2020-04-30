package store

import (
	"context"
	"errors"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type interfaceService interface {
	InsertOne(models.Service) (primitive.ObjectID, error)
	DeleteOne(bson.M) (int64, error)
	DeleteMany(bson.M) (int64, error)
	Find(bson.M) ([]models.Service, error)
	FindOne(bson.M) (models.Service, error)
	FindOneAndUpdate(bson.M, bson.M) (models.Service, error)
}

type service struct {
	db *mongo.Database
}

func (s *service) InsertOne(service models.Service) (primitive.ObjectID, error) {
	collection := s.db.Collection(CollectionServices)

	result, err := collection.InsertOne(context.TODO(), service)
	if err != nil {
		return primitive.NewObjectID(), errors.New("an error occured while saving service to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (s *service) DeleteOne(filter bson.M) (int64, error) {
	collection := s.db.Collection(CollectionServices)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (s *service) DeleteMany(filter bson.M) (int64, error) {
	collection := s.db.Collection(CollectionServices)
	res, err := collection.DeleteMany(context.TODO(), filter)

	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (s *service) Find(filter bson.M) ([]models.Service, error) {
	var services []models.Service

	collection := s.db.Collection(CollectionServices)

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	for cur.Next(context.TODO()) {
		var service models.Service

		err = cur.Decode(&service)
		if err != nil {
			return nil, err
		}

		services = append(services, service)
	}

	return services, nil
}

func (s *service) FindOne(filter bson.M) (models.Service, error) {
	var service models.Service

	collection := s.db.Collection(CollectionServices)

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

func (s *service) FindOneAndUpdate(filter, update bson.M) (models.Service, error) {
	collection := s.db.Collection(CollectionServices)

	res := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
	)
	if res.Err() != nil {
		return models.Service{}, nil
	}

	var service models.Service

	err := res.Decode(&service)
	if err != nil {
		return models.Service{}, err
	}

	return service, nil
}
