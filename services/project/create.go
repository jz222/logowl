package project

import (
	"context"
	"errors"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/organization"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(project models.Project) (primitive.ObjectID, error) {
	timestamp := time.Now()
	project.CreatedAt = timestamp
	project.UpdatedAt = timestamp

	if !project.Validate() {
		return primitive.ObjectID{}, errors.New("the provided project data is invalid")
	}

	organizationExists, err := organization.CheckPresence(bson.M{"_id": project.OrganizationID})
	if err != nil {
		return primitive.ObjectID{}, err
	}
	if !organizationExists {
		return primitive.ObjectID{}, errors.New("the provided organization does not exist")
	}

	collection := mongodb.GetClient().Collection("projects")

	result, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving project to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
