package project

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
)

func Find(filter bson.M) ([]models.Project, error) {
	var projects []models.Project

	projects = []models.Project{}

	collection := mongodb.GetClient().Collection("projects")

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return []models.Project{}, err
	}

	for cur.Next(context.TODO()) {
		var project models.Project

		err = cur.Decode(&project)
		if err != nil {
			return []models.Project{}, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}
