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

var (
	courseMongoNamespace = db.Namespace{Database: "praromvik", Collection: "courses"}
)

func (c *Course) Get() error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	result, err := mongoDB.GetDocument(bson.D{{Key: "uuid", Value: c.UUID}})
	if err != nil {
		return err
	}
	if err := result.Decode(&c); err != nil {
		return err
	}
	return nil
}

func (c *Course) Delete() error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	if err := mongoDB.DeleteDocument(bson.D{{Key: "uuid", Value: c.UUID}}); err != nil {
		return err
	}
	return nil
}

func (c *Course) List() ([]Course, error) {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	cursor, err := mongoDB.ListDocument()
	if err != nil {
		return nil, err
	}
	var list []Course
	if err = cursor.All(context.Background(), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (c *Course) ValidateCourseIDUniqueness() (int, error) {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	result, err := mongoDB.GetDocument(bson.D{{Key: "courseId", Value: c.CourseId}})
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return http.StatusBadRequest, fmt.Errorf("course id '%s' already exists", c.CourseId)
	}

	return http.StatusOK, nil
}

func (c *Course) AddCourseDataToDB() error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	err := mongoDB.AddDocument(c)
	return err
}

func GetCourseUUID(courseID string) (string, error) {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{courseMongoNamespace}}
	result, err := mongoDB.GetDocument(bson.D{{Key: "courseId", Value: courseID}})
	if err != nil {
		return "", err
	}
	course := &Course{}
	if err := result.Decode(course); err != nil {
		return "", err
	}
	return course.UUID, nil
}
