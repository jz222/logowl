package user

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckPresence(filter bson.M) (bool, error) {
	collection := mongodb.GetClient().Collection("users")
	count, err := collection.CountDocuments(context.TODO(), filter)

	return count > 0, err
}
