package event

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"go.mongodb.org/mongo-driver/bson"
)

func UpdateError(filter, update bson.M) error {
	collection := mongodb.GetClient().Collection("errors")

	res := collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
