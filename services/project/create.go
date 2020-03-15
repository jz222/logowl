package project

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(project models.Project) (primitive.ObjectID, error) {
	timestamp := time.Now()
	project.CreatedAt = timestamp
	project.UpdatedAt = timestamp

	if !project.Validate() {
		return primitive.ObjectID{}, errors.New("the provided project data is invalid")
	}

	collection := mongodb.GetClient().Collection("projects")

	result, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		log.Println("Failed to save new project to database with error:", err.Error())
		return primitive.ObjectID{}, errors.New("an error occured while saving project to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
