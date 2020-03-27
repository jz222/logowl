package organization

import (
	"context"
	"fmt"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Delete(organizationID primitive.ObjectID) error {
	collection := mongodb.GetClient().Collection("services")

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

		collection := mongodb.GetClient().Collection("services")
		_, err := collection.DeleteMany(context.TODO(), bson.M{"_id": bson.M{"$in": allServiceIDs}})
		if err != nil {
			fmt.Println("services")
		}
		c <- err
	}()

	go func() {
		if len(allTickets) == 0 {
			c <- nil
			return
		}

		collection := mongodb.GetClient().Collection("errors")
		_, err := collection.DeleteMany(context.TODO(), bson.M{"ticket": bson.M{"$in": allTickets}})
		if err != nil {
			fmt.Println("errors")
		}
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection("organizations")
		_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": organizationID})
		if err != nil {
			fmt.Println("orga")
		}
		c <- err
	}()

	go func() {
		collection := mongodb.GetClient().Collection("users")
		_, err := collection.DeleteMany(context.TODO(), bson.M{"organizationId": organizationID})
		if err != nil {
			fmt.Println("users")
		}
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
