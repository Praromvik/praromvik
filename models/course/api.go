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
	"google.golang.org/genproto/googleapis/type/date"
)

type Course struct {
	UUID               string    `json:"uuid" bson:"uuid"`
	CourseId           string    `json:"courseId" bson:"courseId"`
	Title              string    `json:"title" bson:"title"`
	Description        string    `json:"description" bson:"description"`
	Instructors        []string  `json:"instructors" bson:"instructors"`
	Moderators         []string  `json:"moderators" bson:"moderators"`
	StartDate          date.Date `json:"startDate" bson:"startDate"`
	EndDate            date.Date `json:"endDate" bson:"endDate"`
	Duration           int       `json:"duration" bson:"duration"` // Duration in week
	EnrollmentStudents []string  `json:"enrollmentStudents" bson:"enrollmentStudents"`
	Lessons            []string  `json:"lessons" bson:"lessons"`
	Price              string    `json:"price" bson:"price"`
}

type Lesson struct {
	UUID     string   `json:"uuid" bson:"uuid"`
	LessonId string   `json:"lessonId" bson:"lessonId"`
	Title    string   `json:"title" bson:"title"`
	Content  []string `json:"content" bson:"content"`
}

type Content struct {
	UUID string `json:"uuid" bson:"uuid"`
	Name string `json:"Name" bson:"name"`
	Type string `json:"type" bson:"type"` // video, resource, quiz, lab
}

type Quiz struct {
	Questions []string `json:"questions" bson:"questions"`
}

type Question struct {
	UUID           string   `json:"uuid" bson:"uuid"`
	Name           string   `json:"Name" bson:"Name"`
	Options        []string `json:"options" bson:"options"`
	CorrectOptions []string `json:"correctOptions" bson:"correctOptions"`
}
