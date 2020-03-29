package user

import (
	"context"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Update(filter, update bson.M) error {
	collection := mongodb.GetClient().Collection(mongodb.Users)

	newPassword, ok := update["password"]
	if ok {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword.(string)), 12)
		if err != nil {
			return err
		}

		update["password"] = string(hash)
	}

	update["updatedAt"] = time.Now()

	res := collection.FindOneAndUpdate(context.TODO(), filter, bson.M{"$set": update})
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
