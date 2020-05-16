package store

import (
	"context"
	"errors"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interfaceOrganization interface {
	CheckPresence(bson.M) (bool, error)
	InsertOne(models.Organization) (primitive.ObjectID, error)
	DeleteOne(bson.M) (int64, error)
	FindOne(bson.M) (models.Organization, error)
	FindOneAndUpdate(filter, update bson.M) (models.Organization, error)
}

type organization struct {
	db *mongo.Database
}

func (o *organization) CheckPresence(filter bson.M) (bool, error) {
	collection := o.db.Collection(CollectionOrganizations)
	count, err := collection.CountDocuments(context.TODO(), filter, options.Count().SetLimit(1))

	return count > 0, err
}

func (o *organization) InsertOne(organization models.Organization) (primitive.ObjectID, error) {
	collection := o.db.Collection(CollectionOrganizations)

	result, err := collection.InsertOne(context.TODO(), organization)
	if err != nil {
		return primitive.NilObjectID, errors.New("an error occured while saving organization to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (o *organization) DeleteOne(filter bson.M) (int64, error) {
	collection := o.db.Collection(CollectionOrganizations)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (o *organization) FindOne(filter bson.M) (models.Organization, error) {
	var organization models.Organization

	collection := o.db.Collection(CollectionOrganizations)

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

func (o *organization) FindOneAndUpdate(filter, update bson.M) (models.Organization, error) {
	collection := o.db.Collection(CollectionOrganizations)

	res := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
		options.MergeFindOneAndUpdateOptions().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return models.Organization{}, res.Err()
	}

	var organization models.Organization

	err := res.Decode(&organization)
	if err != nil {
		return models.Organization{}, err
	}

	return organization, nil
}
