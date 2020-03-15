package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/jz222/loggy/keys"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

// InitiateDatabase creates a new connection to MongoDB that can then
// be retrieved by using the GetClient function.
func InitiateDatabase() {
	if db != nil {
		return
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(
		ctx,
		options.Client().ApplyURI(keys.GetKeys().MONGO_URI),
	)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("❌", err.Error())
	}

	log.Println("✅ Connection to MongoDB established")

	db = client.Database(keys.GetKeys().MONGO_DB_NAME)
}

// GetClient returns a MongoDB instance.
func GetClient() *mongo.Database {
	if db != nil {
		return db
	}

	InitiateDatabase()

	return db
}
