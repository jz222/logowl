package mongodb

import (
	"context"
	"log"
	"time"

	"github.com/jz222/loggy/keys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var db *mongo.Database

const (
	Errors        = "errors"
	Organizations = "organizations"
	Services      = "services"
	Users         = "users"
)

func initiateDatabase() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(keys.GetKeys().MONGO_URI))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB with error: ", err.Error())
	}

	db = client.Database(keys.GetKeys().MONGO_DB_NAME)

	collection := db.Collection(Errors)
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.M{"fingerprint": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"updatedAt": -1},
			Options: nil,
		},
	}
	collection.Indexes().CreateMany(ctx, indexModels)

	collection = db.Collection(Users)
	indexModels = []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	collection.Indexes().CreateMany(ctx, indexModels)

	log.Println("✅ Connection to MongoDB established")
}

// GetClient returns a MongoDB instance.
func GetClient() *mongo.Database {
	if db != nil {
		return db
	}

	initiateDatabase()

	return db
}
