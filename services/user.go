package services

import (
	"context"
	"errors"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

type user struct {
	DB *mongo.Database
}

func (u *user) FetchAllInformation(filter bson.M) (models.User, error) {
	ctx := context.Background()
	collection := u.DB.Collection(mongodb.Users)

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

	cur, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return models.User{}, err
	}
	defer cur.Close(ctx)

	var doc models.User

	cur.Next(ctx)
	cur.Decode(&doc)

	return doc, nil
}

func (u *user) CheckPresence(filter bson.M) (bool, error) {
	collection := u.DB.Collection(mongodb.Users)
	count, err := collection.CountDocuments(context.TODO(), filter, options.Count().SetLimit(1))

	return count > 0, err
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

	collection := u.DB.Collection(mongodb.Users)

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving user to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (u *user) Delete(filter bson.M) (int64, error) {
	collection := u.DB.Collection(mongodb.Users)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (u *user) FindOne(filter bson.M) (models.User, error) {
	var user models.User

	collection := u.DB.Collection(mongodb.Users)

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

func (u *user) Invite(userData models.User) (models.User, error) {
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

	collection := u.DB.Collection(mongodb.Users)

	result, err := collection.InsertOne(context.TODO(), userData)
	if err != nil {
		return models.User{}, err
	}

	userData.ID = result.InsertedID.(primitive.ObjectID)
	userData.Password = ""

	return userData, nil
}

func (u *user) Update(filter, update bson.M) error {
	collection := u.DB.Collection(mongodb.Users)

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

func GetUserService(db *mongo.Database) user {
	return user{
		DB: db,
	}
}
