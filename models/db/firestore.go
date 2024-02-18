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
	"cloud.google.com/go/firestore"

	"context"
	"fmt"

	"github.com/praromvik/praromvik/models/db/client"
)

type Firestore struct{}

func (_ Firestore) GetDocument(collName string, id string) (*firestore.DocumentSnapshot, error) {
	dSnap, err := client.Firestore.Collection(collName).Doc(id).Get(context.Background())
	if err != nil {
		return nil, err
	}
	return dSnap, nil
}

func (_ Firestore) AddDocument(collName string, id string, data interface{}) error {
	_, err := client.Firestore.Collection(collName).Doc(id).Create(context.Background(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Inserted document with _id: %v\n", id)
	return nil
}

func (_ Firestore) UpdateDocument(collName string, id string, data interface{}) error {
	_, err := client.Firestore.Collection(collName).Doc(id).Set(context.Background(), data)
	if err != nil {
		return err
	}

	fmt.Printf("Update Document with id %s\n", id)
	return nil
}
