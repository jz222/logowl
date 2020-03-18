package organization

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func CheckPresence() (bool, error) {
	collection := mongodb.GetClient().Collection("organizations")
	count, err := collection.CountDocuments(context.TODO(), bson.M{})

	return count > 0, err
}
