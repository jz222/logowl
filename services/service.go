package services

import (
	"errors"
	"time"

	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/store"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
)

type InterfaceService interface {
	CheckPresence(bson.M) (bool, error)
	Create(models.Service) (models.Service, error)
	Delete(bson.M) (int64, error)
	Find(bson.M) (*[]models.Service, error)
	FindOne(bson.M) (*models.Service, error)
}

type service struct {
	DB store.InterfaceStore
}

func (s *service) CheckPresence(filter bson.M) (bool, error) {
	return s.DB.Service().CheckPresence(filter)
}

func (s *service) Create(service models.Service) (models.Service, error) {
	timestamp := time.Now()
	service.CreatedAt = timestamp
	service.UpdatedAt = timestamp

	if !service.Validate() {
		return models.Service{}, errors.New("the provided service data is invalid")
	}

	organizationExists, err := s.DB.Organization().CheckPresence(bson.M{"_id": service.OrganizationID})
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

	result, err := s.DB.Service().InsertOne(service)
	if err != nil {
		return models.Service{}, errors.New("an error occured while saving service to database")
	}

	service.ID = result

	return service, nil
}

func (s *service) Delete(filter bson.M) (int64, error) {
	return s.DB.Service().DeleteOne(filter)
}

func (s *service) Find(filter bson.M) (*[]models.Service, error) {
	return s.DB.Service().Find(filter)
}

func (s *service) FindOne(filter bson.M) (*models.Service, error) {
	return s.DB.Service().FindOne(filter)
}

func GetServiceService(db store.InterfaceStore) service {
	return service{
		DB: db,
	}
}
