package organization

import (
	"context"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Delete(organizationID primitive.ObjectID) error {
	collection := mongodb.GetClient().Collection(mongodb.Services)

	cur, err := collection.Find(context.TODO(), bson.M{"organizationId": organizationID})
	if err != nil {
		return err
	}

	var allServiceIDs []primitive.ObjectID
	var allTickets []string

	for cur.Next(context.TODO()) {
		var service models.Service

		err := cur.Decode(&service)
		if err != nil {
			return nil
		}

		allServiceIDs = append(allServiceIDs, service.ID)
		allTickets = append(allTickets, service.Ticket)
	}

	c := make(chan error, 4)

	go func() {
		if len(allServiceIDs) == 0 {
			c <- nil
			return
		}

		collection := mongodb.GetClient().Collection(mongodb.Services)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": allServiceIDs}})
		c <- err
	}()

	go func() {
		if len(allTickets) == 0 {
			c <- nil
			return
		}

		collection := mongodb.GetClient().Collection(mongodb.Errors)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"ticket": bson.M{"$in": allTickets}})
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection(mongodb.Organizations)
		_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": organizationID})
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection(mongodb.Users)
		_, err := collection.DeleteMany(context.TODO(), bson.M{"organizationId": organizationID})
		c <- err
	}()

	var failed error

	for i := 0; i < 4; i++ {
		err := <-c

		if err != nil {
			failed = err
		}
	}

	return failed
}
