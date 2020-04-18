package services

import (
	"time"

	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterfaceEvent interface {
	GetError(bson.M, primitive.ObjectID) (models.Error, error)
	GetErrors(string, int64) ([]models.Error, error)
	DeleteError(bson.M) (int64, error)
	UpdateError(bson.M, bson.M) error
}

type event struct {
	Store store.InterfaceStore
}

func (e *event) DeleteError(filter bson.M) (int64, error) {
	return e.Store.Error().DeleteOne(filter)
}

func (e *event) GetError(filter bson.M, viewer primitive.ObjectID) (models.Error, error) {
	return e.Store.Error().FindOneAndUpdate(filter, bson.M{"$addToSet": bson.M{"seenBy": viewer}}, true)
}

func (e *event) GetErrors(ticket string, page int64) ([]models.Error, error) {
	return e.Store.Error().FindPaged(bson.M{"ticket": ticket}, page)
}

func (e *event) UpdateError(filter, update bson.M) error {
	update["updatedAt"] = time.Now()

	_, err := e.Store.Error().FindOneAndUpdate(filter, bson.M{"$set": update}, false)
	if err != nil {
		return err
	}

	return nil
}

func GetEventService(store store.InterfaceStore) event {
	return event{store}
}
