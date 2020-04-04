package services

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InterfaceLogging interface {
	SaveError(models.Error)
}

type logging struct {
	DB *mongo.Database
}

func (l *logging) SaveError(errorEvent models.Error) {
	serviceService := GetServiceService(l.DB)

	serviceExists, err := serviceService.CheckPresence(bson.M{"ticket": errorEvent.Ticket})
	if err != nil {
		log.Println("Failed to verify service with error:", err.Error())
	}

	if !serviceExists || err != nil {
		return
	}

	hash := md5.Sum([]byte(errorEvent.Message + errorEvent.Stacktrace + errorEvent.Ticket))

	_, convertedTimestamp, err := utils.FormatTimestamp(errorEvent.Timestamp)
	if err != nil {
		log.Println("failed to convert timestamp:", errorEvent.Timestamp)
		return
	}

	timestamp := time.Now()

	errorEvent.Fingerprint = hex.EncodeToString(hash[:])
	errorEvent.Evolution = map[string]int{convertedTimestamp: 1}
	errorEvent.Count = 1
	errorEvent.LastSeen = errorEvent.Timestamp
	errorEvent.CreatedAt = timestamp
	errorEvent.UpdatedAt = timestamp

	collection := l.DB.Collection(mongodb.Errors)

	_, err = collection.InsertOne(context.TODO(), errorEvent)
	if err == nil {
		return
	}

	key := fmt.Sprintf("%s.%s", "evolution", convertedTimestamp)

	collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"fingerprint": errorEvent.Fingerprint},
		bson.M{
			"$inc": bson.M{"count": 1, key: 1},
			"$set": bson.M{"lastSeen": errorEvent.Timestamp, "updatedAt": timestamp},
		},
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
	)
}

func GetLoggingService(db *mongo.Database) logging {
	return logging{
		DB: db,
	}
}
