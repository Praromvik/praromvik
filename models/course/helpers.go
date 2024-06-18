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
	"github.com/praromvik/praromvik/models/db"
)

type Document interface {
	GetID() interface{}
	GetNamespace() db.Namespace
}

func (c *Course) GetNamespace() db.Namespace {
	return db.Namespace{
		Database:   "praromvik",
		Collection: "courses",
	}
}

func (l *Lesson) GetNamespace() db.Namespace {
	return db.Namespace{
		Database:   l.CourseRef,
		Collection: "lessons",
	}
}

func (c *Content) GetNamespace() db.Namespace {
	return db.Namespace{
		Database:   c.CourseRef,
		Collection: "contents",
	}
}

func (l *Lesson) GetID() interface{} {
	return l.LessonID
}

func (c *Course) GetID() interface{} {
	return c.CourseId
}

func (c *Content) GetID() interface{} {
	return c.ContentID
}

// Another Approach to set Value to pointer.
//func (c *Course) SetDocument(document interface{}) error {
//	doc, ok := document.(*Course)
//	if !ok {
//		return fmt.Errorf("document is not of type %s", reflect.TypeOf(c))
//	}
//	*c = *doc
//	return nil
//}
