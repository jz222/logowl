package services

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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InterfaceService interface {
	CheckPresence(bson.M) (bool, error)
	Create(models.Service) (models.Service, error)
	Delete(bson.M) (int64, error)
	Find(bson.M) ([]models.Service, error)
	FindOne(bson.M) (models.Service, error)
}

type service struct {
	DB *mongo.Database
}

func (s *service) CheckPresence(filter bson.M) (bool, error) {
	collection := s.DB.Collection(mongodb.Services)
	count, err := collection.CountDocuments(context.TODO(), filter, options.Count().SetLimit(1))

	return count > 0, err
}

func (s *service) Create(service models.Service) (models.Service, error) {
	timestamp := time.Now()
	service.CreatedAt = timestamp
	service.UpdatedAt = timestamp

	if !service.Validate() {
		return models.Service{}, errors.New("the provided service data is invalid")
	}

	organizationExists, err := organization.CheckPresence(bson.M{"_id": service.OrganizationID})
	if err != nil {
		return models.Service{}, err
	}
	if !organizationExists {
		return models.Service{}, errors.New("the provided organization does not exist")
	}

	ticket, err := utils.GenerateTicket()
	if err != nil {
		return models.Service{}, err
	}

	service.Ticket = ticket

	collection := s.DB.Collection(mongodb.Services)

	result, err := collection.InsertOne(context.TODO(), service)
	if err != nil {
		return models.Service{}, errors.New("an error occured while saving service to database")
	}

	service.ID = result.InsertedID.(primitive.ObjectID)

	return service, nil
}

func (s *service) Delete(filter bson.M) (int64, error) {
	collection := s.DB.Collection(mongodb.Services)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (s *service) Find(filter bson.M) ([]models.Service, error) {
	var services []models.Service

	services = []models.Service{}

	collection := s.DB.Collection(mongodb.Services)

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return []models.Service{}, err
	}

	for cur.Next(context.TODO()) {
		var service models.Service

		err = cur.Decode(&service)
		if err != nil {
			return []models.Service{}, err
		}

		services = append(services, service)
	}

	return services, nil
}

func (s *service) FindOne(filter bson.M) (models.Service, error) {
	var service models.Service

	collection := s.DB.Collection(mongodb.Services)

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

func GetServiceService(db *mongo.Database) service {
	return service{
		DB: db,
	}
}
