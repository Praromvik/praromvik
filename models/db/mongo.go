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

	"github.com/praromvik/praromvik/models/db/client"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo struct {
	Namespaces Namespaces
}

func (m Mongo) GetDocument(filter interface{}) (*mongo.SingleResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.FindOne(context.TODO(), filter), nil
}

func (m Mongo) AddDocument(data interface{}) (*mongo.InsertOneResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.InsertOne(context.TODO(), data)
}

func (m Mongo) UpdateDocument(filter interface{}, newData interface{}) (*mongo.UpdateResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.ReplaceOne(context.TODO(), filter, newData)
}

func (m Mongo) DeleteDocument(filter interface{}) (*mongo.DeleteResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.DeleteOne(context.TODO(), filter)
}

func (m Mongo) ListDocuments(filter interface{}) (*mongo.Cursor, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.Find(context.Background(), filter)
}

func (m Mongo) CountDocuments(filter interface{}) (int64, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return 0, err
	}
	return collection.CountDocuments(context.Background(), filter)
}

func (m Mongo) BulkWrite(models []mongo.WriteModel) (*mongo.BulkWriteResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	opts := options.BulkWrite().SetOrdered(false)
	return collection.BulkWrite(context.Background(), models, opts)
}

func (m Mongo) FindDistinct(field string, filter interface{}) ([]interface{}, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.Distinct(context.Background(), field, filter)
}

func (m Mongo) Aggregate(pipeline interface{}) (*mongo.Cursor, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	return collection.Aggregate(context.Background(), pipeline)
}

func (m Mongo) getDBCollection() (*mongo.Collection, error) {
	if len(m.Namespaces) == 0 {
		return nil, fmt.Errorf("empty database name")
	}
	db, collection := &mongo.Database{}, &mongo.Collection{}
	for _, val := range m.Namespaces {
		db = client.Mongo.Database(val.Database)
		collection = db.Collection(val.Collection)
	}
	return collection, nil
}

func (m Mongo) Sync(query string, id string, bsonName string, elementId string) error {
	collection, err := m.getDBCollection()
	if err != nil {
		return err
	}
	// Check if the field is null (unset), and if so, initialize it as an empty array
	var doc bson.M
	err = collection.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		return fmt.Errorf("failed to find document: %v", err)
	}

	if doc[bsonName] == nil {
		// Initialize the field as an empty array
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": id},
			bson.M{"$set": bson.M{bsonName: []string{elementId}}},
		)
		if err != nil {
			return fmt.Errorf("failed to initialize field: %v", err)
		}
	} else {
		// Field is not null, proceed with pushing the element
		_, err = collection.UpdateOne(
			context.TODO(),
			bson.M{"_id": id},
			bson.M{query: bson.M{bsonName: elementId}},
		)
		if err != nil {
			return fmt.Errorf("failed to sync document: %v", err)
		}
	}
	return err
}
