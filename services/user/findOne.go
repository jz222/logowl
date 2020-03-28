package user

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FindOne(filter bson.M) (models.User, error) {
	var user models.User

	collection := mongodb.GetClient().Collection(mongodb.Users)

	queryResult := collection.FindOne(context.TODO(), filter)
	if queryResult.Err() != nil {
		return models.User{}, queryResult.Err()
	}

	err := queryResult.Decode(&user)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
