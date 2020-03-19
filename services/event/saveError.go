package event

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/jz222/loggy/libs/mongodb"
	"github.com/jz222/loggy/models"
	"github.com/jz222/loggy/services/project"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SaveError(errorLog models.Error) {
	projectExists, err := project.CheckPresence(bson.M{"ticket": errorLog.Ticket})
	if err != nil {
		log.Println("Failed to verify project with error:", err.Error())
	}

	if !projectExists {
		return
	}

	hash := md5.Sum([]byte(errorLog.Message + errorLog.Stacktrace))
	errorLog.Fingerprint = hex.EncodeToString(hash[:])

	collection := mongodb.GetClient().Collection("errors")
	newDoc := collection.FindOneAndUpdate(
		context.TODO(),
		bson.M{"fingerprint": errorLog.Fingerprint},
		bson.M{"message": "updated"},
		options.MergeFindOneAndUpdateOptions().SetUpsert(true),
	)

	fmt.Println(newDoc)
}
