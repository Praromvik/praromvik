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

package course

import (
	"context"
	"fmt"
	"go/types"
	"net/http"
	"reflect"

	"github.com/praromvik/praromvik/models/db"
	"go.mongodb.org/mongo-driver/bson"
)

func Get(document Document) (any, error) {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	result, err := mongoDB.GetDocument(bson.D{{Key: "_id", Value: document.GetID()}})
	if err != nil {
		return nil, err
	}

	// Create a new instance of the same type as document to decode into
	documentType := reflect.TypeOf(document).Elem()
	newDoc := reflect.New(documentType).Interface()

	if err := result.Decode(newDoc); err != nil {
		return nil, err
	}
	return newDoc, nil

	// ***Another Approach***
	//// We can't decode the document hence its only interface, therefore we need another method like 'SetDocument' to original struct
	//// Set the decoded document back to the original document
	//if err := document.SetDocument(newDoc); err != nil {
	//	return err
	//}
}

func Create(document Document) error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	_, err := mongoDB.AddDocument(document)
	return err
}

func Delete(document Document) error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	if _, err := mongoDB.DeleteDocument(bson.D{{Key: "_id", Value: document.GetID()}}); err != nil {
		return err
	}
	return nil
}

func Update(document Document) error {
	filter := bson.D{{Key: "_id", Value: document.GetID()}}
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	result, err := mongoDB.GetDocument(filter)
	if err != nil {
		return err
	}

	// Create a new instance of the same type as document to decode into
	existingDocType := reflect.TypeOf(document).Elem()
	existingDoc := reflect.New(existingDocType).Interface()

	// Decode the existing document
	if err := result.Decode(existingDoc); err != nil {
		return err
	}

	// Merge the updated fields into the existing document
	db.MergeStruct(existingDoc, document)

	// Update the document in the database
	_, err = mongoDB.UpdateDocument(filter, existingDoc)
	return err
}

func List(document Document) (any, error) {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	cursor, err := mongoDB.ListDocuments(types.Interface{})
	if err != nil {
		return nil, err
	}

	// Create a slice of the correct type
	documentType := reflect.TypeOf(document).Elem()
	sliceType := reflect.SliceOf(documentType)
	documents := reflect.New(sliceType).Interface()

	// Decode all documents into the slice
	if err := cursor.All(context.Background(), documents); err != nil {
		return nil, err
	}
	return documents, nil
}

func Sync(document Document, query string, id string, bsonName string, elementID string) error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	if err := mongoDB.Sync(query, id, bsonName, elementID); err != nil {
		return err
	}
	return nil
}

func ValidateNameUniqueness(document Document) (int, error) {
	filter := bson.D{{Key: "_id", Value: document.GetID()}}
	mongoDB := db.Mongo{Namespaces: []db.Namespace{document.GetNamespace()}}
	count, err := mongoDB.CountDocuments(filter)
	if err != nil {
		return http.StatusBadRequest, err
	}

	if count != 0 {
		return http.StatusBadRequest, fmt.Errorf("course id '%s' already exists", document.GetID())
	}
	return http.StatusOK, nil
}
