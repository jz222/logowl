package user

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Delete(id primitive.ObjectID) (int64, error) {
	collection := mongodb.GetClient().Collection(mongodb.Users)

	res, err := collection.DeleteOne(context.TODO(), bson.M{"_id": id})
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}
