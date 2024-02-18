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
	"fmt"

	"github.com/praromvik/praromvik/models/db/client"

	"go.mongodb.org/mongo-driver/mongo"

	"reflect"
)

func getDBCollection(dbAndCollList []string) (*mongo.Collection, error) {
	if len(dbAndCollList) == 0 {
		return nil, fmt.Errorf("empty database name")
	} else if len(dbAndCollList)%2 == 1 {
		return nil, fmt.Errorf("invalid number of DB and collection names provided, must be even")
	}
	db, collection := &mongo.Database{}, &mongo.Collection{}
	for i, val := range dbAndCollList {
		if i%2 == 0 {
			db = client.Mongo.Database(val)
		} else {
			collection = db.Collection(val)
		}
	}
	return collection, nil
}

func MergeStruct(oldStruct interface{}, newStruct interface{}) {
	oldValue, newValue := reflect.ValueOf(oldStruct).Elem(), reflect.ValueOf(newStruct)
	for i := 0; i < oldValue.NumField(); i++ {
		oldFieldValue, newFieldValue := oldValue.Field(i), newValue.Field(i)
		if oldFieldValue.Kind() == reflect.Struct {
			MergeStruct(oldFieldValue.Addr().Interface(), newFieldValue.Interface())
		} else if !reflect.DeepEqual(newFieldValue.Interface(), reflect.Zero(newFieldValue.Type()).Interface()) {
			oldFieldValue.Set(newFieldValue)
		}
	}
}
