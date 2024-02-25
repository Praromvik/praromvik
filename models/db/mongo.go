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
	"go.mongodb.org/mongo-driver/mongo"
)

type Mongo struct {
	Namespaces Namespaces
}

func (m Mongo) GetDocument(filter interface{}) (*mongo.SingleResult, error) {
	collection, err := m.getDBCollection()
	if err != nil {
		return nil, err
	}
	result := collection.FindOne(context.TODO(), filter)
	return result, nil
}

func (m Mongo) AddDocument(data interface{}) error {
	collection, err := m.getDBCollection()
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

func (m Mongo) UpdateDocument(filter interface{}, newData interface{}) error {
	collection, err := m.getDBCollection()
	if err != nil {
		return err
	}
	result, err := collection.ReplaceOne(context.TODO(), filter, newData)
	if err != nil {
		return fmt.Errorf("failed to update document: %v", err)
	}
	fmt.Printf("Update Document. Result. MatchedCount:"+
		" %d, UpsertedCount: %d, ModifiedCount: %d.\n", result.MatchedCount, result.UpsertedCount, result.ModifiedCount)
	return nil
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
