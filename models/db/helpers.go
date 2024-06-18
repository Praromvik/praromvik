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

package db

import (
	"reflect"
)

type Namespaces []Namespace

type Namespace struct {
	Database   string
	Collection string
}

// Here, You must have to provide Structure kind instead of ptr.

func MergeStruct(oldStruct interface{}, newStruct interface{}) {
	// We're taking pointer as parameter, but we need struct value instead of pointer
	// therefore [Value.Elem()] gives us struct instead of pointer.
	oldValue, newValue := reflect.ValueOf(oldStruct).Elem(), reflect.ValueOf(newStruct).Elem()
	for i := 0; i < oldValue.NumField(); i++ {
		oldFieldValue, newFieldValue := oldValue.Field(i), newValue.Field(i)
		if oldFieldValue.CanSet() {
			if oldFieldValue.Kind() == reflect.Struct {
				// At this point we've to pass pointer as parameter, so we convert this struct to pointer.
				MergeStruct(oldFieldValue.Addr().Interface(), newFieldValue.Addr().Interface())
			} else if !reflect.DeepEqual(newFieldValue.Interface(), reflect.Zero(newFieldValue.Type()).Interface()) {
				// [Value.Set] can set value only if field is addressable. we can ensure it by calling [Value.CanSet]
				oldFieldValue.Set(newFieldValue)
			}
		}
	}
}
