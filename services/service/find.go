package service

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Find(filter bson.M) ([]models.Service, error) {
	var services []models.Service

	services = []models.Service{}

	collection := mongodb.GetClient().Collection(mongodb.Services)

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return []models.Service{}, err
	}

	for cur.Next(context.TODO()) {
		var service models.Service

		err = cur.Decode(&service)
		if err != nil {
			return []models.Service{}, err
		}

		services = append(services, service)
	}

	return services, nil
}
