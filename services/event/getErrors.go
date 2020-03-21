package event

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetErrors(ticket, pointer string) (*[]models.Error, error) {
	collection := mongodb.GetClient().Collection("errors")

	var filter bson.M

	if pointer == "" {
		filter = bson.M{"ticket": ticket}
	} else {
		id, err := primitive.ObjectIDFromHex(pointer)
		if err != nil {
			return nil, err
		}

		filter = bson.M{"ticket": ticket, "_id": bson.M{"$lt": id}}
	}

	cur, err := collection.Find(
		context.TODO(),
		filter,
		options.MergeFindOptions().SetSort(bson.M{"updatedAt": -1}),
		options.MergeFindOptions().SetLimit(5),
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

	return &errorEvents, nil
}

func Populate() {
	collection := mongodb.GetClient().Collection("errors")
	var manyEvents []interface{}

	for i := 0; i <= 200000; i++ {
		event := models.Error{
			Message:   "some error " + strconv.Itoa(i),
			Ticket:    "testticket",
			CreatedAt: time.Now().Add(time.Duration(i) * time.Millisecond),
			UpdatedAt: time.Now().Add(time.Duration(i) * time.Millisecond),
		}

		hash := md5.Sum([]byte(event.Message + event.Stacktrace))
		event.Fingerprint = hex.EncodeToString(hash[:])

		manyEvents = append(manyEvents, event)
	}

	fmt.Println("uploading")

	_, err := collection.InsertMany(context.TODO(), manyEvents)
	if err != nil {
		fmt.Println(err.Error())
	}
}
