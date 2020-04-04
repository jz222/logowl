package services

import (
	"errors"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/store"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type InterfaceUser interface {
	FetchAllInformation(bson.M) (*models.User, error)
	CheckPresence(bson.M) (bool, error)
	Create(models.User) (primitive.ObjectID, error)
	Delete(bson.M) (int64, error)
	FindOne(bson.M) (*models.User, error)
	Invite(models.User) (*models.User, error)
	Update(bson.M, bson.M) error
}

type user struct {
	DB store.InterfaceStore
}

func (u *user) FetchAllInformation(filter bson.M) (*models.User, error) {
	pipeline := []bson.M{
		bson.M{
			"$match": filter,
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "organizationId",
				"from":         mongodb.Organizations,
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
				"from":         mongodb.Services,
				"foreignField": "organizationId",
				"as":           "services",
			},
		},
		bson.M{
			"$lookup": bson.M{
				"localField":   "organizationId",
				"from":         mongodb.Users,
				"foreignField": "organizationId",
				"as":           "team",
			},
		},
	}

	user, err := u.DB.User().Aggregate(pipeline)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *user) CheckPresence(filter bson.M) (bool, error) {
	return u.DB.User().CheckPresence(filter)
}

func (u *user) Create(user models.User) (primitive.ObjectID, error) {
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
	user.IsVerified = true

	result, err := u.DB.User().InsertOne(user)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving user to database")
	}

	return result, nil
}

func (u *user) Delete(filter bson.M) (int64, error) {
	return u.DB.User().DeleteOne(filter)
}

func (u *user) FindOne(filter bson.M) (*models.User, error) {
	return u.DB.User().FindOne(filter)
}

func (u *user) Invite(userData models.User) (*models.User, error) {
	timestamp := time.Now()
	userData.CreatedAt = timestamp
	userData.UpdatedAt = timestamp

	randomString, err := utils.GenerateRandomString(12)
	if err != nil {
		return nil, err
	}

	userData.Password = randomString

	if !userData.Validate() {
		return nil, errors.New("the provided user data is invalid")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(userData.Password), 12)
	if err != nil {
		return nil, err
	}

	userData.Password = string(hash)

	inviteCode, err := utils.GenerateRandomString(20)
	if err != nil {
		return nil, err
	}

	userData.InviteCode = inviteCode
	userData.IsVerified = false

	result, err := u.DB.User().InsertOne(userData)
	if err != nil {
		return nil, err
	}

	userData.ID = result
	userData.Password = ""

	return &userData, nil
}

func (u *user) Update(filter, update bson.M) error {
	newPassword, ok := update["password"]
	if ok {
		hash, err := bcrypt.GenerateFromPassword([]byte(newPassword.(string)), 12)
		if err != nil {
			return err
		}

		update["password"] = string(hash)
	}

	update["updatedAt"] = time.Now()

	err := u.DB.User().FindOneAndUpdate(filter, bson.M{"$set": update})
	if err != nil {
		return err
	}

	return nil
}

func GetUserService(db store.InterfaceStore) user {
	return user{
		DB: db,
	}
}
