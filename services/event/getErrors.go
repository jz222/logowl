package event

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetErrors(ticket string, page int64) (*[]models.Error, error) {
	collection := mongodb.GetClient().Collection(mongodb.Errors)

	cur, err := collection.Find(
		context.TODO(),
		bson.M{"ticket": ticket},
		options.MergeFindOptions().SetSort(bson.M{"updatedAt": -1}),
		options.MergeFindOptions().SetSkip(page*10),
		options.MergeFindOptions().SetLimit(10),
	)
	if err != nil {
		return nil, err
	}

	var errorEvents []models.Error

	for cur.Next(context.TODO()) {
		var errorEvent models.Error

		err := cur.Decode(&errorEvent)
		if err == nil {
			errorEvents = append(errorEvents, errorEvent)
		}
	}

	if errorEvents == nil {
		errorEvents = []models.Error{}
	}

	return &errorEvents, nil
}
