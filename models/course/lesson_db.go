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

func (l *Lesson) getDbNamespace() []db.Namespace {
	return []db.Namespace{{Database: l.CourseRef, Collection: "lessons"}}
}

func (l *Lesson) ValidateNameUniqueness() (int, error) {
	mongoDB := db.Mongo{Namespaces: l.getDbNamespace()}
	result, err := mongoDB.GetDocument(bson.D{{Key: "_id", Value: l.LessonID}})
	if err != nil {
		return http.StatusBadRequest, err
	}
	if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return http.StatusBadRequest, fmt.Errorf("lesson id '%s' already exists", l.LessonID)
	}
	return http.StatusOK, nil
}

func (l *Lesson) Create() error {
	mongoDB := db.Mongo{Namespaces: l.getDbNamespace()}
	err := mongoDB.AddDocument(l)
	return err
}

func (l *Lesson) Get() error {
	mongoDB := db.Mongo{Namespaces: l.getDbNamespace()}
	result, err := mongoDB.GetDocument(bson.D{{Key: "_id", Value: l.LessonID}})
	if err != nil {
		return err
	}
	if err := result.Decode(&l); err != nil {
		return err
	}
	return nil
}

func (l *Lesson) List() ([]Lesson, error) {
	mongoDB := db.Mongo{Namespaces: l.getDbNamespace()}
	cursor, err := mongoDB.ListDocument()
	if err != nil {
		return nil, err
	}
	var list []Lesson
	if err = cursor.All(context.Background(), &list); err != nil {
		return nil, err
	}
	return list, nil
}

func (l *Lesson) Delete() error {
	if err := l.Get(); err != nil {
		return err
	}
	if l.Contents != nil {
		return errors.New("for deletion, the contents of the lesson must not be nil")
	}

	mongoDB := db.Mongo{Namespaces: l.getDbNamespace()}
	if err := mongoDB.DeleteDocument(bson.D{{Key: "_id", Value: l.LessonID}}); err != nil {
		return err
	}

	return nil
}
