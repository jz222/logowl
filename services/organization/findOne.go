package organization

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FindOne(filter bson.M) (models.Organization, error) {
	var organization models.Organization

	collection := mongodb.GetClient().Collection("organizations")

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
