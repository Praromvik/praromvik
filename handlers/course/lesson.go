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
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/praromvik/praromvik/models/course"
	perror "github.com/praromvik/praromvik/pkg/error"

	"github.com/go-chi/chi/v5"
)

type Lesson struct {
	*course.Lesson
}

func (l *Lesson) Create(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&l.Lesson); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
		return
	}
	l.CourseRef = chi.URLParam(r, "courseRef")
	errCode, err := course.ValidateNameUniqueness(l.Lesson)
	if err != nil {
		perror.HandleError(w, errCode, "", err)
		return
	}
	if err := course.Create(l.Lesson); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "failed to create course lesson data into database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (l *Lesson) Get(w http.ResponseWriter, r *http.Request) {
	l.Lesson = &course.Lesson{
		CourseRef: chi.URLParam(r, "courseRef"),
		LessonID:  chi.URLParam(r, "id"),
	}
	// Fetch document from database
	document, err := course.Get(l.Lesson)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course lesson", err)
		return
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&course.Lesson{})
	if !isTypeValid(document, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the document to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document.(*course.Lesson)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (l *Lesson) List(w http.ResponseWriter, r *http.Request) {
	l.Lesson = &course.Lesson{
		CourseRef: chi.URLParam(r, "courseRef"),
		LessonID:  chi.URLParam(r, "id"),
	}
	documents, err := course.List(l.Lesson)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course lesson list.", err)
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&[]course.Lesson{})
	if !isTypeValid(documents, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the documents to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents.(*[]course.Lesson)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
}

func (l *Lesson) Delete(w http.ResponseWriter, r *http.Request) {
	l.Lesson = &course.Lesson{
		CourseRef: chi.URLParam(r, "courseRef"),
		LessonID:  chi.URLParam(r, "id"),
	}
	if err := course.Delete(l.Lesson); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on deleting course lesson.", err)
	}
	w.WriteHeader(http.StatusOK)
}
