package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/jz222/loggy/internal/keys"
	"github.com/jz222/loggy/internal/models"
	"github.com/jz222/loggy/internal/store"
	"github.com/jz222/loggy/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type InterfaceUser interface {
	FetchAllInformation(bson.M) (models.User, error)
	CheckPresence(bson.M) (bool, error)
	Create(models.User) (primitive.ObjectID, error)
	Delete(bson.M) (int64, error)
	FindOne(bson.M) (models.User, error)
	Invite(models.User) (models.User, error)
	Update(bson.M, bson.M) error
}

type User struct {
	Store   store.InterfaceStore
	Request InterfaceRequest
}

func (u *User) FetchAllInformation(filter bson.M) (models.User, error) {
	pipeline := []bson.M{
		bson.M{
			"$match": filter,
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "organizationId",
				"from":         store.CollectionOrganizations,
				"foreignField": "_id",
				"as":           "organization",
			},
		},
		bson.M{
			"$unwind": bson.M{
				"path":                       "$organization",
				"preserveNullAndEmptyArrays": true,
			},
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "organizationId",
				"from":         store.CollectionServices,
				"foreignField": "organizationId",
				"as":           "services",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "organizationId",
				"from":         store.CollectionUsers,
				"foreignField": "organizationId",
				"as":           "team",
			},
		},
	}

	user, err := u.Store.User().Aggregate(pipeline)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (u *User) CheckPresence(filter bson.M) (bool, error) {
	return u.Store.User().CheckPresence(filter)
}

func (u *User) Create(user models.User) (primitive.ObjectID, error) {
	timestamp := time.Now()
	user.CreatedAt = timestamp
	user.UpdatedAt = timestamp

	if !user.Validate() {
		return primitive.NilObjectID, errors.New("the provided user data is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
	if err != nil {
		return primitive.NilObjectID, err
	}

	user.Password = string(hash)
	user.IsVerified = true

	result, err := u.Store.User().InsertOne(user)
	if err != nil {
		return primitive.NilObjectID, errors.New("an error occured while saving user to database")
	}

	return result, nil
}

func (u *User) Delete(filter bson.M) (int64, error) {
	return u.Store.User().DeleteOne(filter)
}

func (u *User) FindOne(filter bson.M) (models.User, error) {
	return u.Store.User().FindOne(filter)
}

func (u *User) Invite(userData models.User) (models.User, error) {
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
		return models.User{}, err
	}

	userData.InviteCode = inviteCode
	userData.IsVerified = false

	result, err := u.Store.User().InsertOne(userData)
	if err != nil {
		return models.User{}, err
	}

	userData.ID = result
	userData.Password = ""

	emailData := map[string]interface{}{
		"FirstName": userData.FirstName,
		"URL":       fmt.Sprintf("%s/auth/signup?code=%s", keys.GetKeys().CLIENT_URL, userData.InviteCode),
	}

	go u.Request.SendEmail(userData.Email, "invitation", emailData)

	return userData, nil
}

func (u *User) Update(filter, update bson.M) error {
	newPassword, ok := update["password"]
	if ok {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword.(string)), 12)
		if err != nil {
			return err
		}

		update["password"] = string(hash)
	}

	update["updatedAt"] = time.Now()

	err := u.Store.User().FindOneAndUpdate(filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}

func GetUserService(store store.InterfaceStore) User {
	return User{store, &Request{}}
}
