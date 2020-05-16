package store

import (
	"context"
	"errors"

	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type interfacePasswordResetTokens interface {
	InsertOne(models.PasswordResetToken) (primitive.ObjectID, error)
	FindOneAndUpdate(bson.M, bson.M) (models.PasswordResetToken, error)
}

type passwordResetTokens struct {
	db *mongo.Database
}

func (p *passwordResetTokens) InsertOne(passwordResetToken models.PasswordResetToken) (primitive.ObjectID, error) {
	collection := p.db.Collection(CollectionPasswordResetTokens)

	result, err := collection.InsertOne(context.TODO(), passwordResetToken)
	if err != nil {
		return primitive.NilObjectID, errors.New("an error occured while saving password reset token to the the database")
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (p *passwordResetTokens) FindOneAndUpdate(filter, update bson.M) (models.PasswordResetToken, error) {
	collection := p.db.Collection(CollectionPasswordResetTokens)

	res := collection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
		options.MergeFindOneAndUpdateOptions().SetReturnDocument(options.After),
	)
	if res.Err() != nil {
		return models.PasswordResetToken{}, res.Err()
	}

	var passwordResetToken models.PasswordResetToken

	err := res.Decode(&passwordResetToken)
	if err != nil {
		return models.PasswordResetToken{}, err
	}

	return passwordResetToken, nil
}
