package organization

import (
	"context"
	"errors"
	"regexp"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Create(organization models.Organization) (primitive.ObjectID, error) {
	timestamp := time.Now()
	organization.CreatedAt = timestamp
	organization.UpdatedAt = timestamp

	if !organization.Validate() {
		return primitive.ObjectID{}, errors.New("the provided organization data is invalid")
	}

	regex := regexp.MustCompile(`\s+`)
	organization.Identifier = regex.ReplaceAllString(organization.Name, "")

	collection := mongodb.GetClient().Collection(mongodb.Organizations)

	result, err := collection.InsertOne(context.TODO(), organization)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving organization to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
