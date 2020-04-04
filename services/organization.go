package services

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InterfaceOrganization interface {
	CheckPresence(bson.M) (bool, error)
	Create(models.Organization) (primitive.ObjectID, error)
	Delete(primitive.ObjectID) error
	FindOne(bson.M) (models.Organization, error)
}

type organization struct {
	DB *mongo.Database
}

func (o *organization) CheckPresence(filter bson.M) (bool, error) {
	collection := o.DB.Collection(mongodb.Organizations)
	count, err := collection.CountDocuments(context.TODO(), filter, options.Count().SetLimit(1))

	return count > 0, err
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

	collection := o.DB.Collection(mongodb.Organizations)

	result, err := collection.InsertOne(context.TODO(), organization)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving organization to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *organization) Delete(organizationID primitive.ObjectID) error {
	collection := o.DB.Collection(mongodb.Services)

	cur, err := collection.Find(context.TODO(), bson.M{"organizationId": organizationID})
	if err != nil {
		return err
	}

	var allServiceIDs []primitive.ObjectID
	var allTickets []string

	for cur.Next(context.TODO()) {
		var service models.Service

		err := cur.Decode(&service)
		if err != nil {
			return nil
		}

		allServiceIDs = append(allServiceIDs, service.ID)
		allTickets = append(allTickets, service.Ticket)
	}

	c := make(chan error, 4)

	go func() {
		if len(allServiceIDs) == 0 {
			c <- nil
			return
		}

		collection := mongodb.GetClient().Collection(mongodb.Services)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": allServiceIDs}})
		c <- err
	}()

	go func() {
		if len(allTickets) == 0 {
			c <- nil
			return
		}

		collection := mongodb.GetClient().Collection(mongodb.Errors)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"ticket": bson.M{"$in": allTickets}})
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection(mongodb.Organizations)
		_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": organizationID})
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection(mongodb.Users)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"organizationId": organizationID})
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

func (o *organization) FindOne(filter bson.M) (models.Organization, error) {
	var organization models.Organization

	collection := o.DB.Collection(mongodb.Organizations)

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return models.Organization{}, queryResult.Err()
	}

	err := queryResult.Decode(&organization)
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil
}

func GetOrganizationService(db *mongo.Database) organization {
	return organization{
		DB: db,
	}
}
