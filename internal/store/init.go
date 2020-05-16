package store

import (
	"context"
	"log"
	"time"

	"github.com/jz222/loggy/internal/keys"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	CollectionAnalytics           = "analytics"
	CollectionErrors              = "errors"
	CollectionOrganizations       = "organizations"
	CollectionServices            = "services"
	CollectionUsers               = "users"
	CollectionPasswordResetTokens = "passwordResetTokens"
)

type InterfaceStore interface {
	Connect()
	Disconnect()
	User() interfaceUser
	Service() interfaceService
	Organization() interfaceOrganization
	Error() interfaceErrorEvent
	Analytics() interfaceAnalytics
	PasswordResetTokens() interfacePasswordResetTokens
}

type store struct {
	db *mongo.Database
}

func (s *store) Connect() {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(keys.GetKeys().MONGO_URI))
	if err != nil {
		log.Fatal(err.Error())
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("❌ Failed to connect to MongoDB with error: ", err.Error())
	}

	s.db = client.Database(keys.GetKeys().MONGO_DB_NAME)

	collection := s.db.Collection(CollectionErrors)
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

	collection = s.db.Collection(CollectionUsers)
	indexModels = []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	collection.Indexes().CreateMany(ctx, indexModels)

	collection = s.db.Collection(CollectionAnalytics)
	indexModels = []mongo.IndexModel{
		{
			Keys:    bson.M{"ticket": 1, "month": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	collection.Indexes().CreateMany(ctx, indexModels)

	collection = s.db.Collection(CollectionPasswordResetTokens)
	indexModels = []mongo.IndexModel{
		{
			Keys:    bson.M{"createdAt": 1},
			Options: options.Index().SetExpireAfterSeconds(60),
		},
	}
	collection.Indexes().CreateMany(ctx, indexModels)

	log.Println("✅ Connection to MongoDB established")
}

func (s *store) Disconnect() {
	s.db.Client().Disconnect(context.TODO())

	log.Println("✅ Successfully disconnected from MongoDB")
}

func (s *store) User() interfaceUser {
	return &user{s.db}
}

func (s *store) Service() interfaceService {
	return &service{s.db}
}

func (s *store) Organization() interfaceOrganization {
	return &organization{s.db}
}

func (s *store) Error() interfaceErrorEvent {
	return &errorEvent{s.db}
}

func (s *store) Analytics() interfaceAnalytics {
	return &analytics{s.db}
}

func (s *store) PasswordResetTokens() interfacePasswordResetTokens {
	return &passwordResetTokens{s.db}
}

func GetStore() InterfaceStore {
	return &store{}
}
