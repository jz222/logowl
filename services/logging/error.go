package logging

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"log"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/service"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveError(errorEvent models.Error) {
	serviceExists, err := service.CheckPresence(bson.M{"ticket": errorEvent.Ticket})
	if err != nil {
		log.Println("Failed to verify service with error:", err.Error())
	}

	if !serviceExists || err != nil {
		return
	}

	hash := md5.Sum([]byte(errorEvent.Message + errorEvent.Stacktrace))

	errorEvent.Fingerprint = hex.EncodeToString(hash[:])
	errorEvent.CreatedAt = time.Now()
	errorEvent.UpdatedAt = time.Now()

	collection := mongodb.GetClient().Collection("errors")

	_, err = collection.InsertOne(context.TODO(), errorEvent)
	if err == nil {
		return
	}

	collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"fingerprint": errorEvent.Fingerprint},
		bson.M{"$inc": bson.M{"count": 1}, "$set": bson.M{"updatedAt": time.Now()}},
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
	)
}
