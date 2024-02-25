/*
MIT License

Copyright (c) 2024 Praromvik

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package client

import (
	"cloud.google.com/go/firestore"
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

var (
	Firestore *firestore.Client
	Mongo     *mongo.Client
	Redis     *redis.Client
)

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
