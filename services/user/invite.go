package user

import (
	"context"
	"errors"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Invite(userData models.User) (models.User, error) {
	timestamp := time.Now()
	userData.CreatedAt = timestamp
	userData.UpdatedAt = timestamp

	randomString, err := utils.GenerateRandomString(12)
	if err != nil {
		return models.User{}, err
	}

	userData.Password = randomString

	if !userData.Validate() {
		return models.User{}, errors.New("the provided user data is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)
	if err != nil {
		return models.User{}, err
	}

	userData.Password = string(hash)

	inviteCode, err := utils.GenerateRandomString(20)
	if err != nil {
		return models.User{}, nil
	}

	userData.InviteCode = inviteCode
	userData.IsVerified = false

	collection := mongodb.GetClient().Collection(mongodb.Users)

	result, err := collection.InsertOne(context.TODO(), userData)
	if err != nil {
		return models.User{}, err
	}

	userData.ID = result.InsertedID.(primitive.ObjectID)
	userData.Password = ""

	return userData, nil
}
