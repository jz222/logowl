package services

import (
	"errors"
	"time"

	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type InterfaceService interface {
	CheckPresence(bson.M) (bool, error)
	Create(models.Service) (models.Service, error)
	Delete(bson.M) (int64, error)
	Find(bson.M) ([]models.Service, error)
	FindOne(bson.M) (models.Service, error)
	FindOneAndUpdate(bson.M, bson.M) (models.Service, error)
}

type service struct {
	Store store.InterfaceStore
}

func (s *service) CheckPresence(filter bson.M) (bool, error) {
	return s.Store.Service().CheckPresence(filter)
}

func (s *service) Create(service models.Service) (models.Service, error) {
	timestamp := time.Now()
	service.CreatedAt = timestamp
	service.UpdatedAt = timestamp

	if !service.Validate() {
		return models.Service{}, errors.New("the provided service data is invalid")
	}

	organizationExists, err := s.Store.Organization().CheckPresence(bson.M{"_id": service.OrganizationID})
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

	result, err := s.Store.Service().InsertOne(service)
	if err != nil {
		return models.Service{}, errors.New("an error occured while saving service to database")
	}

	service.ID = result

	return service, nil
}

func (s *service) Delete(filter bson.M) (int64, error) {
	return s.Store.Service().DeleteOne(filter)
}

func (s *service) Find(filter bson.M) ([]models.Service, error) {
	return s.Store.Service().Find(filter)
}

func (s *service) FindOne(filter bson.M) (models.Service, error) {
	return s.Store.Service().FindOne(filter)
}

func (s *service) FindOneAndUpdate(filter, update bson.M) (models.Service, error) {
	update["updatedAt"] = time.Now()

	return s.Store.Service().FindOneAndUpdate(filter, bson.M{"$set": update})
}

func GetServiceService(store store.InterfaceStore) service {
	return service{store}
}
