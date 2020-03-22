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
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetErrors(ticket string, page int64) (*[]models.Error, error) {
	collection := mongodb.GetClient().Collection("errors")

	cur, err := collection.Find(
		context.TODO(),
		bson.M{"ticket": ticket},
		options.MergeFindOptions().SetSort(bson.M{"updatedAt": -1}),
		options.MergeFindOptions().SetSkip(page*5),
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

	if errorEvents == nil {
		errorEvents = []models.Error{}
	}

	return &errorEvents, nil
}

func Populate() {
	collection := mongodb.GetClient().Collection("errors")
	var manyEvents []interface{}

	for i := 1; i <= 200000; i++ {

		ticket := "236F3D82B655CA37083E6561C982D282FA1B03D9573B2FB1A2"
		if i <= 100000 {
			ticket = "8280EFF0B5AEB9957EAFD7FEF78D655C3F5E14C4F120B42CB4"
		}

		event := models.Error{
			Message:   "some error " + strconv.Itoa(i),
			Ticket:    ticket,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Millisecond),
			UpdatedAt: time.Now().Add(time.Duration(i) * time.Millisecond),
		}

		hash := md5.Sum([]byte(event.Message + event.Stacktrace))
		event.Fingerprint = hex.EncodeToString(hash[:])

		fmt.Println(event.ID)

		manyEvents = append(manyEvents, event)
	}

	fmt.Println("uploading")

	_, err := collection.InsertMany(context.TODO(), manyEvents)
	if err != nil {
		fmt.Println(err.Error())
	}
}
