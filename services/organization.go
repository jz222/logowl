package services

import (
	"errors"
	"regexp"
	"time"

	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/store"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type InterfaceOrganization interface {
	CheckPresence(bson.M) (bool, error)
	Create(models.Organization) (primitive.ObjectID, error)
	Delete(primitive.ObjectID) error
	FindOne(bson.M) (*models.Organization, error)
}

type organization struct {
	DB store.InterfaceStore
}

func (o *organization) CheckPresence(filter bson.M) (bool, error) {
	return o.DB.Organization().CheckPresence(filter)
}

func (o *organization) Create(organization models.Organization) (primitive.ObjectID, error) {
	timestamp := time.Now()
	organization.CreatedAt = timestamp
	organization.UpdatedAt = timestamp

	if !organization.Validate() {
		return primitive.ObjectID{}, errors.New("the provided organization data is invalid")
	}

	regex := regexp.MustCompile(`\s+`)
	organization.Identifier = regex.ReplaceAllString(organization.Name, "")

	return o.DB.Organization().InsertOne(organization)
}

func (o *organization) Delete(organizationID primitive.ObjectID) error {
	allServices, err := o.DB.Service().Find(bson.M{"organizationId": organizationID})
	if err != nil {
		return err
	}

	var allServiceIDs []primitive.ObjectID
	var allTickets []string

	for _, service := range *allServices {
		allServiceIDs = append(allServiceIDs, service.ID)
		allTickets = append(allTickets, service.Ticket)
	}

	c := make(chan error, 4)

	go func() {
		if len(allServiceIDs) == 0 {
			c <- nil
			return
		}

		_, err := o.DB.Service().DeleteMany(bson.M{"_id": bson.M{"$in": allServiceIDs}})
		c <- err
	}()

	go func() {
		if len(allTickets) == 0 {
			c <- nil
			return
		}

		_, err := o.DB.Error().DeleteMany(bson.M{"ticket": bson.M{"$in": allTickets}})
		c <- err
	}()

	go func() {
		_, err := o.DB.Organization().DeleteOne(bson.M{"_id": organizationID})
		c <- err
	}()

	go func() {
		_, err := o.DB.User().DeleteMany(bson.M{"organizationId": organizationID})
		c <- err
	}()

	var failed error

	for i := 0; i < 4; i++ {
		err := <-c

		if err != nil {
			failed = err
		}
	}

	return failed
}

func (o *organization) FindOne(filter bson.M) (*models.Organization, error) {
	return o.DB.Organization().FindOne(filter)
}

func GetOrganizationService(db store.InterfaceStore) organization {
	return organization{
		DB: db,
	}
}
