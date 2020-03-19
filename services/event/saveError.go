package event

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/project"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveError(errorEvent models.Error) {
	projectExists, err := project.CheckPresence(bson.M{"ticket": errorEvent.Ticket})
	if err != nil {
		log.Println("Failed to verify project with error:", err.Error())
	}

	if !projectExists {
		return
	}

	hash := md5.Sum([]byte(errorEvent.Message + errorEvent.Stacktrace))
	errorEvent.Fingerprint = hex.EncodeToString(hash[:])

	collection := mongodb.GetClient().Collection("errors")

	_, err = collection.InsertOne(context.TODO(), errorEvent)
	if err == nil {
		return
	}

	collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"fingerprint": errorEvent.Fingerprint},
		bson.M{"$inc": bson.M{"count": 1}},
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
	)
}
