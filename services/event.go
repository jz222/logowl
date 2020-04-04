package services

import (
	"time"

	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/store"
	"go.mongodb.org/mongo-driver/bson"
)

type InterfaceEvent interface {
	GetError(bson.M) (*models.Error, error)
	GetErrors(string, int64) (*[]models.Error, error)
	DeleteError(bson.M) (int64, error)
	UpdateError(bson.M, bson.M) error
}

type event struct {
	DB store.InterfaceStore
}

func (e *event) DeleteError(filter bson.M) (int64, error) {
	return e.DB.Error().DeleteOne(filter)
}

func (e *event) GetError(filter bson.M) (*models.Error, error) {
	return e.DB.Error().FindOne(filter)
}

func (e *event) GetErrors(ticket string, page int64) (*[]models.Error, error) {
	return e.DB.Error().FindPaged(bson.M{"ticket": ticket}, page)
}

func (e *event) UpdateError(filter, update bson.M) error {
	update["updatedAt"] = time.Now()

	err := e.DB.Error().FindOneAndUpdate(filter, bson.M{"$set": update}, false)
	if err != nil {
		return err
	}

	return nil
}

func GetEventService(db store.InterfaceStore) event {
	return event{
		DB: db,
	}
}
