/*
MIT License

# Copyright (c) 2024 Praromvik

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
	"github.com/go-chi/chi/v5"
	"net/http"

	"github.com/praromvik/praromvik/handlers/utils"
	"github.com/praromvik/praromvik/models/course"
	"github.com/praromvik/praromvik/pkg/auth"
	perror "github.com/praromvik/praromvik/pkg/error"
)

type Course struct {
	*course.Course
}
type Lesson struct {
	*course.Lesson
}

func (c *Course) Create(w http.ResponseWriter, r *http.Request) {
	if err := json.NewDecoder(r.Body).Decode(&c.Course); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
		return
	}

	errCode, err := c.ValidateNameUniqueness()
	if err != nil {
		perror.HandleError(w, errCode, "", err)
		return
	}

	info, err := auth.GetUserInfoFromSession(r)
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting session", err)
		return
	}
	c.Instructors = append(c.Instructors, *info)
	if err := c.AddCourseDataToDB(); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "failed to course data into database", err)
		return
	}
}

func (c *Course) List(w http.ResponseWriter, r *http.Request) {
	c.Course = &course.Course{}
	courseList, err := c.Course.List()
	if err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course list.", err)
	}
	var mapList []map[string]interface{}
	for i := range courseList {
		val := &courseList[i]
		data, err := utils.RemovedUUID(val)
		if err != nil {
			perror.HandleError(w, http.StatusInternalServerError, "", err)
			return
		}
		mapList = append(mapList, data)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(mapList); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
	}
}

func (c *Course) Get(w http.ResponseWriter, r *http.Request) {
	c.CourseId = chi.URLParam(r, "id")
	if err := c.Course.Get(); err != nil {
		perror.HandleError(w, http.StatusBadRequest, "Error on getting course.", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(c); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "Error on encoding JSON response", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Course) Update(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an Course by Name")
}

//func (c *Course) Delete(w http.ResponseWriter, r *http.Request) {
//	var ok bool
//	c.Course = &course.Course{}
//	ctx := r.Context()
//	c.UUID, ok = ctx.Value("uuid").(string)
//	if !ok {
//		perror.HandleError(w, http.StatusUnprocessableEntity, "",
//			fmt.Errorf("failed to extract uuid from url"))
//		return
//	}
//	if err := c.Course.Delete(); err != nil {
//		perror.HandleError(w, http.StatusBadRequest, "Error on getting course.", err)
//	}
//	w.WriteHeader(http.StatusOK)
//}
