package store

import (
	"context"
	"errors"

	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interfaceUser interface {
	InsertOne(models.User) (primitive.ObjectID, error)
	Aggregate([]bson.M) (models.User, error)
	CheckPresence(bson.M) (bool, error)
	DeleteOne(bson.M) (int64, error)
	DeleteMany(bson.M) (int64, error)
	FindOne(bson.M) (models.User, error)
	FindOneAndUpdate(bson.M, bson.M) error
}

type user struct {
	db *mongo.Database
}

func (u *user) InsertOne(user models.User) (primitive.ObjectID, error) {
	collection := u.db.Collection(CollectionUsers)

	result, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		return primitive.ObjectID{}, errors.New("an error occured while saving user to database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (u *user) Aggregate(pipeline []bson.M) (models.User, error) {
	ctx := context.TODO()
	collection := u.db.Collection(CollectionUsers)

	cur, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return models.User{}, err
	}
	defer cur.Close(ctx)

	var user models.User

	cur.Next(ctx)
	cur.Decode(&user)

	return user, nil
}

func (u *user) CheckPresence(filter bson.M) (bool, error) {
	collection := u.db.Collection(CollectionUsers)
	count, err := collection.CountDocuments(context.TODO(), filter, options.Count().SetLimit(1))

	return count > 0, err
}

func (u *user) DeleteOne(filter bson.M) (int64, error) {
	collection := u.db.Collection(CollectionUsers)

	res, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (u *user) DeleteMany(filter bson.M) (int64, error) {
	collection := u.db.Collection(CollectionUsers)

	res, err := collection.DeleteMany(context.TODO(), filter)
	if err != nil {
		return 0, err
	}

	return res.DeletedCount, nil
}

func (u *user) FindOne(filter bson.M) (models.User, error) {
	var user models.User

	collection := u.db.Collection(CollectionUsers)

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

func (u *user) FindOneAndUpdate(filter, update bson.M) error {
	collection := u.db.Collection(CollectionUsers)

	res := collection.FindOneAndUpdate(context.TODO(), filter, update)
	if res.Err() != nil {
		return res.Err()
	}

	return nil
}
