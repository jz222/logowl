package user

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Create(user models.User) (primitive.ObjectID, error) {
	timestamp := time.Now()
	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp

	if !user.Validate() {
		return primitive.ObjectID{}, errors.New("the provided user data is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	user.Password = string(hash)

	collection := mongodb.GetClient().Collection("users")

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Println("Failed to save new user to database with error:", err.Error())
		return primitive.ObjectID{}, errors.New("an error occured while saving user to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}
