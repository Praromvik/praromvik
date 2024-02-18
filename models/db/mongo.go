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

package db

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct{}

func (_ Mongo) GetDocument(dbAndCollList []string, filter interface{}) (*mongo.SingleResult, error) {
	collection, err := getDBCollection(dbAndCollList)
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(context.TODO(), filter)
	return result, nil
}

func (_ Mongo) AddDocument(dbAndCollList []string, data interface{}) error {
	collection, err := getDBCollection(dbAndCollList)
	if err != nil {
		return err
	}
	result, err := collection.InsertOne(context.TODO(), data)
	if err != nil {
		return err
	}
	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
	return nil
}

func (_ Mongo) UpdateDocument(dbAndCollList []string, filter interface{}, newDocument interface{}) error {
	collection, err := getDBCollection(dbAndCollList)
	if err != nil {
		return err
	}
	result, err := collection.ReplaceOne(context.TODO(), filter, newDocument)
	if err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}
	fmt.Printf("Update Document. Result. MatchedCount:"+
		" %d, UpsertedCount: %d, ModifiedCount: %d.\n", result.MatchedCount, result.UpsertedCount, result.ModifiedCount)
	return nil
}
