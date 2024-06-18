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

type Content struct {
	*course.Content
}

func (c *Content) Create(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&c.Content); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
		return
	}
	c.CourseRef = chi.URLParam(r, "courseRef")
	errCode, err := course.ValidateNameUniqueness(c.Content)
	if err != nil {
		perror.HandleError(w, errCode, "", err)
		return
	}
	if err := course.Create(c.Content); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "failed to create course content data into database", err)
		return
	}

	if err := course.Sync(&course.Lesson{
		CourseRef: c.Content.CourseRef,
	}, "$push", c.LessonRef, "contents", c.ContentID); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on syncing content to lesson", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Content) Get(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
		ContentID: chi.URLParam(r, "id"),
	}
	// Fetch document from database
	document, err := course.Get(c.Content)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course lesson", err)
		return
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&course.Content{})
	if !isTypeValid(document, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the document to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document.(*course.Content)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Content) List(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
		ContentID: chi.URLParam(r, "id"),
	}
	documents, err := course.List(c.Content)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course content list.", err)
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&[]course.Content{})
	if !isTypeValid(documents, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the documents to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents.(*[]course.Content)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
}

func (c *Content) Delete(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
		ContentID: chi.URLParam(r, "id"),
	}
	// Fetch document from database
	document, err := course.Get(c.Content)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course lesson", err)
		return
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&course.Content{})
	if !isTypeValid(document, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}
	c.Content.LessonRef = document.(*course.Content).LessonRef

	if err := course.Delete(c.Content); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on deleting content.", err)
	}

	if err := course.Sync(&course.Lesson{
		CourseRef: c.Content.CourseRef,
	}, "$pull", c.Content.LessonRef, "contents", c.ContentID); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on syncing content to lesson", err)
	}

	w.WriteHeader(http.StatusOK)
}
