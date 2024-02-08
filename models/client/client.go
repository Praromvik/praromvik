package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var Firestore *firestore.Client
var Mongo *mongo.Client
var Redis *redis.Client

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error getting env, %v", err)
	}

	ctx := context.Background()
	var err error
	Firestore, err = ConnectToFireStore(ctx)
	if err != nil {
		log.Fatalf("failed to get firestore client: %v", err)
	}

	Mongo, err = ConnectToMongoDB(ctx)
	if err != nil {
		log.Fatalf("failed to get MongoDB client: %v", err)
	}

	Redis, err = ConnectToRedis()
	if err != nil {
		log.Fatalf("failed to get redis client: %v", err)
	}
}
