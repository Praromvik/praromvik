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
	errCode, err := c.ValidateNameUniqueness()
	if err != nil {
		perror.HandleError(w, errCode, "", err)
		return
	}
	if err := c.Content.Create(); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "failed to create course data into database", err)
		return
	}
}

func (c *Content) Get(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
		ContentID: chi.URLParam(r, "id"),
	}
	fmt.Println("-----------Inside Get------------")
	fmt.Println("")
	if err := c.Content.Get(); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting content.", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Content) List(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
	}
	list, err := c.Content.List()
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting content list.", err)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(list); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
	}
}

func (c *Content) Delete(w http.ResponseWriter, r *http.Request) {
	c.Content = &course.Content{
		CourseRef: chi.URLParam(r, "courseRef"),
		ContentID: chi.URLParam(r, "id"),
	}
	if err := c.Content.Delete(); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on deleting content.", err)
	}
	w.WriteHeader(http.StatusOK)
}
