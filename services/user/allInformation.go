package user

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func FetchAllInformation(filter bson.M) (models.User, error) {
	ctx := context.Background()
	collection := mongodb.GetClient().Collection(mongodb.Users)

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
