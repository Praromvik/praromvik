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
	"slices"

	"github.com/praromvik/praromvik/models/course"
	"github.com/praromvik/praromvik/pkg/auth"
	perror "github.com/praromvik/praromvik/pkg/error"

	"github.com/go-chi/chi/v5"
)

type Course struct {
	*course.Course
}

func (c *Course) Create(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&c.Course); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
		return
	}
	errCode, err := course.ValidateNameUniqueness(c.Course)
	if err != nil {
		perror.HandleError(w, errCode, "", err)
		return
	}
	info, err := auth.GetUserInfoFromSession(r)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting session", err)
		return
	}
	if !slices.Contains(c.Instructors, info.Name) {
		c.Instructors = append(c.Instructors, info.Name)
	}
	if err := course.Create(c.Course); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "failed to course data into database", err)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Course) List(w http.ResponseWriter, r *http.Request) {
	// Initialize Course instance
	c.Course = &course.Course{}
	// Fetch documents from the database
	documents, err := course.List(c.Course)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course list.", err)
		return
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&[]course.Course{})
	if !isTypeValid(documents, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the documents to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(documents.(*[]course.Course)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
}

func (c *Course) Get(w http.ResponseWriter, r *http.Request) {
	// Initialize Course instance
	c.Course = &course.Course{}
	c.CourseId = chi.URLParam(r, "id")
	// Fetch document from database
	document, err := course.Get(c.Course)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course.", err)
		return
	}

	// Check if the document is of the expected type
	expectedType := reflect.TypeOf(&course.Course{})
	if !isTypeValid(document, expectedType) {
		perror.HandleError(w, http.StatusBadRequest, "", fmt.Errorf("document is not of type %s", expectedType))
		return
	}

	// Encode the document to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(document.(*course.Course)); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Course) Update(w http.ResponseWriter, r *http.Request) {
	c.Course = &course.Course{}
	c.CourseId = chi.URLParam(r, "id")
	if err := json.NewDecoder(r.Body).Decode(&c.Course); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
		return
	}
	if err := course.Update(c.Course); err != nil {
		fmt.Println(err)
		perror.HandleError(w, http.StatusBadRequest, "Error on updating course", err)
		return
	}
}

func (c *Course) Delete(w http.ResponseWriter, r *http.Request) {
	c.Course = &course.Course{}
	c.CourseId = chi.URLParam(r, "id")
	if err := course.Delete(c.Course); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on deleting course.", err)
	}
	w.WriteHeader(http.StatusOK)
}

func isTypeValid(document interface{}, expectedType reflect.Type) bool {
	docValue := reflect.ValueOf(document)
	return docValue.Kind() == reflect.Ptr && docValue.Type() == expectedType
}
