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

type Course struct {
	CourseId    string   `json:"_id" bson:"_id"`
	Title       string   `json:"title" bson:"title"`
	Description string   `json:"description" bson:"description"`
	Instructors []string `json:"instructors" bson:"instructors"`
	Moderators  []string `json:"moderators" bson:"moderators"`
	StartDate   string   `json:"startDate" bson:"startDate"`
	EndDate     string   `json:"endDate" bson:"endDate"`
	Duration    int      `json:"duration" bson:"duration"` // Duration in week
	Capacity    int      `json:"capacity" bson:"capacity"`
	Students    []string `json:"students" bson:"students"`
	Price       int      `json:"price" bson:"price"`
	Image       []byte   `json:"image" bson:"-"`
}

type Lesson struct {
	LessonID  string   `json:"_id" bson:"_id"`
	CourseRef string   `json:"courseRef" bson:"courseRef"`
	Title     string   `json:"title" bson:"title"`
	Contents  []string `json:"contents" bson:"contents"`
}

type Content struct {
	ContentID string `json:"_id" bson:"_id"`
	CourseRef string `json:"courseRef" bson:"courseRef"`
	LessonRef string `json:"lessonRef" bson:"lessonRef"`
	Title     string `json:"title" bson:"title"`
	Type      string `json:"type" bson:"type"` // video, resource, quiz, lab
	Data      []byte `json:"data" bson:"data"`
}
