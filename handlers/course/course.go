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
	"fmt"
	"net/http"

	"github.com/praromvik/praromvik/models/course"
)

type Course struct {
	*course.Course
}

func (o *Course) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an Course")
}

func (o *Course) List(w http.ResponseWriter, r *http.Request) {
	fmt.Println("List all Courses")
}

func (o *Course) GetByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Get an Course by ID")
}

func (o *Course) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an Course by ID")
}

func (o *Course) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an Course by ID")
}
