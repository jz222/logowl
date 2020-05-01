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
	Create(models.Service) (models.Service, error)
	Delete(bson.M) (int64, error)
	Find(bson.M) ([]models.Service, error)
	FindOne(bson.M) (models.Service, error)
	FindOneAndUpdate(bson.M, bson.M) (models.Service, error)
}

type Service struct {
	Store store.InterfaceStore
}

func (s *Service) Create(service models.Service) (models.Service, error) {
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

	analytics := models.Analytics{
		Ticket:    ticket,
		Data:      map[string]models.AnalyticData{},
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}

	_, err = s.Store.Analytics().InsertOne(analytics)
	if err != nil {
		return models.Service{}, err
	}

	return service, nil
}

func (s *Service) Delete(filter bson.M) (int64, error) {
	service, err := s.Store.Service().FindOne(filter)
	if err != nil {
		return 0, err
	}

	_, err = s.Store.Error().DeleteMany(bson.M{"ticket": service.Ticket})
	if err != nil {
		return 0, err
	}

	_, err = s.Store.Analytics().DeleteMany(bson.M{"ticket": service.Ticket})
	if err != nil {
		return 0, err
	}

	return s.Store.Service().DeleteOne(filter)
}

func (s *Service) Find(filter bson.M) ([]models.Service, error) {
	return s.Store.Service().Find(filter)
}

func (s *Service) FindOne(filter bson.M) (models.Service, error) {
	return s.Store.Service().FindOne(filter)
}

func (s *Service) FindOneAndUpdate(filter, update bson.M) (models.Service, error) {
	update["updatedAt"] = time.Now()

	return s.Store.Service().FindOneAndUpdate(filter, bson.M{"$set": update})
}

func GetServiceService(store store.InterfaceStore) Service {
	return Service{store}
}
