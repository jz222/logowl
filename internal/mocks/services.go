package mocks

import (
	"github.com/jz222/loggy/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LoggingService struct{}

func (l *LoggingService) SaveError(e models.Error) {}

func (l *LoggingService) SaveAnalyticEvent(a models.AnalyticEvent) {}

type UserService struct{}

func (u *UserService) FetchAllInformation(f bson.M) (models.User, error) {
	return models.User{}, nil
}

func (u *UserService) CheckPresence(f bson.M) (bool, error) {
	return false, nil
}

func (u *UserService) Create(user models.User) (primitive.ObjectID, error) {
	return primitive.NilObjectID, nil
}

func (u *UserService) Delete(f bson.M) (int64, error) {
	return 0, nil
}

func (u *UserService) FindOne(f bson.M) (models.User, error) {
	return models.User{}, nil
}

func (u *UserService) Invite(user models.User) (models.User, error) {
	return models.User{}, nil
}

func (u *UserService) Update(f, update bson.M) error {
	return nil
}
