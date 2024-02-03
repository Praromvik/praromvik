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

package api

import (
	"google.golang.org/genproto/googleapis/type/date"
)

type Course struct {
	CourseId           string       `json:"courseId"`
	Title              string       `json:"title"`
	Description        string       `json:"description"`
	Instructors        []Instructor `json:"instructors"`
	StartDate          date.Date    `json:"startDate"`
	EndDate            date.Date    `json:"endDate"`
	Duration           int          `json:"duration"` // Duration in week
	EnrollmentCapacity int          `json:"enrollmentCapacity"`
	EnrollmentStudents []User       `json:"enrollmentStudents"`
	Lessons            []Lesson     `json:"lessons"`
}

type Lesson struct {
	LessonId string `json:"lessonId"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	// resources []resource `json:"resources"`
	Quizzes []Quiz `json:"quizzes"`
}

type Quiz struct {
	QuizId    string     `json:"quizId"`
	Questions []Question `json:"questions"`
}

type Question struct {
	QuestionId    string   `json:"QuestionId"`
	Ques          string   `json:"ques"`
	Options       []string `json:"options"`
	CorrectOption string   `json:"correctOption"`
}
