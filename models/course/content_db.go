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
	"errors"
	"fmt"
	"net/http"

	"github.com/praromvik/praromvik/models/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (c *Content) getDbNamespace() []db.Namespace {
	return []db.Namespace{{Database: c.CourseRef, Collection: "contents"}}
}

func (c *Content) ValidateNameUniqueness() (int, error) {
	mongoDB := db.Mongo{Namespaces: c.getDbNamespace()}
	fmt.Println("Namespace: ", mongoDB.Namespaces)
	result, err := mongoDB.GetDocument(bson.D{{Key: "_id", Value: c.ContentID}})
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return http.StatusBadRequest, fmt.Errorf("content id '%s' already exists", c.ContentID)
	}
	return http.StatusOK, nil
}

func (c *Content) Create() error {
	mongoDB := db.Mongo{Namespaces: c.getDbNamespace()}
	if err := mongoDB.AddDocument(c); err != nil {
		return err
	}

	mongoDB = db.Mongo{Namespaces: []db.Namespace{{Database: c.CourseRef, Collection: "lessons"}}}
	if err := mongoDB.Sync("$push", c.LessonRef, "contents", c.ContentID); err != nil {
		return err
	}
	return nil
}

func (c *Content) Get() error {
	mongoDB := db.Mongo{Namespaces: c.getDbNamespace()}
	result, err := mongoDB.GetDocument(bson.D{{Key: "_id", Value: c.ContentID}})
	if err != nil {
		return err
	}
	if err := result.Decode(&c); err != nil {
		return err
	}
	return nil
}

func (c *Content) List() ([]Content, error) {
	mongoDB := db.Mongo{Namespaces: c.getDbNamespace()}
	cursor, err := mongoDB.ListDocument()
	if err != nil {
		return nil, err
	}
	var list []Content
	if err = cursor.All(context.Background(), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Content) Delete() error {
	if err := c.Get(); err != nil {
		return err
	}
	mongoDB := db.Mongo{Namespaces: c.getDbNamespace()}
	if err := mongoDB.DeleteDocument(bson.D{{Key: "_id", Value: c.ContentID}}); err != nil {
		return err
	}
	mongoDB = db.Mongo{
		Namespaces: []db.Namespace{
			{Database: c.CourseRef, Collection: "lessons"},
		},
	}
	if err := mongoDB.Sync("$pull", c.LessonRef, "contents", c.ContentID); err != nil {
		return err
	}
	return nil
}
