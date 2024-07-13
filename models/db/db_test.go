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
	"testing"
)

func TestMergeStruct(t *testing.T) {
	type embeddedUnexported struct {
		Field3 string
	}
	type EmbeddedExported struct {
		Field3 string
	}
	tests := []struct {
		name      string
		oldStruct interface{}
		newStruct interface{}
		expected  interface{}
	}{
		{
			name: "ReplaceAllFields",
			oldStruct: &struct {
				Name string
				Age  int
			}{"Anisur", 25},
			newStruct: &struct {
				Name string
				Age  int
			}{"Arnob", 28},
			expected: &struct {
				Name string
				Age  int
			}{"Arnob", 28},
		},
		{
			name: "SetEmptyFields",
			oldStruct: &struct {
				Name string
				Age  int
			}{},
			newStruct: &struct {
				Name string
				Age  int
			}{"Arnob", 28},
			expected: &struct {
				Name string
				Age  int
			}{"Arnob", 28},
		},
		{
			name: "IgnoreUnexportedFields",
			oldStruct: &struct {
				Name   string
				Id     int64
				salary int64
			}{"Anisur", 12345, 130000},
			newStruct: &struct {
				Name   string
				Id     int64
				salary int64
			}{Name: "Arnob", Id: 123456},
			expected: &struct {
				Name   string
				Id     int64
				salary int64
			}{"Arnob", 123456, 130000},
		},
		{
			name: "NestedStructUpdate",
			oldStruct: &struct {
				Name string
				Bio  struct {
					Address    string
					PostalCode int64
				}
			}{Name: "Anisur", Bio: struct {
				Address    string
				PostalCode int64
			}{Address: "Dewliabari, Konabari, Gazipur", PostalCode: 1270}},

			newStruct: &struct {
				Name string
				Bio  struct {
					Address    string
					PostalCode int64
				}
			}{Name: "Anisur", Bio: struct {
				Address    string
				PostalCode int64
			}{Address: "Secotr 10, Uttara, Dhaka", PostalCode: 1210}},
			expected: &struct {
				Name string
				Bio  struct {
					Address    string
					PostalCode int64
				}
			}{Name: "Anisur", Bio: struct {
				Address    string
				PostalCode int64
			}{Address: "Secotr 10, Uttara, Dhaka", PostalCode: 1210}},
		},
		{
			name: "PartialUpdate",
			oldStruct: &struct {
				Name  string
				Age   int
				Email string
			}{"Anisur", 25, "anisur@example.com"},
			newStruct: &struct {
				Name  string
				Age   int
				Email string
			}{Age: 30},
			expected: &struct {
				Name  string
				Age   int
				Email string
			}{"Anisur", 30, "anisur@example.com"},
		},
		{
			name: "NestedStructPartialUpdate",
			oldStruct: &struct {
				Name string
				Bio  struct {
					Address string
					Phone   string
				}
			}{Name: "Anisur", Bio: struct {
				Address string
				Phone   string
			}{Address: "Old Address", Phone: "1234567890"}},

			newStruct: &struct {
				Name string
				Bio  struct {
					Address string
					Phone   string
				}
			}{Bio: struct {
				Address string
				Phone   string
			}{Address: "New Address"}},
			expected: &struct {
				Name string
				Bio  struct {
					Address string
					Phone   string
				}
			}{Name: "Anisur", Bio: struct {
				Address string
				Phone   string
			}{Address: "New Address", Phone: "1234567890"}},
		},
		{
			name: "EmptyNewStruct",
			oldStruct: &struct {
				Name  string
				Age   int
				Email string
			}{"Anisur", 25, "anisur@example.com"},
			newStruct: &struct {
				Name  string
				Age   int
				Email string
			}{},
			expected: &struct {
				Name  string
				Age   int
				Email string
			}{"Anisur", 25, "anisur@example.com"},
		},
		{
			name: "EmbeddedStructUnexported",
			oldStruct: &struct {
				Field1 string
				Field2 string
				embeddedUnexported
			}{Field1: "Field1", Field2: "Field2", embeddedUnexported: embeddedUnexported{Field3: "Field3"}},
			newStruct: &struct {
				Field1 string
				Field2 string
				embeddedUnexported
			}{Field1: "NewField1", Field2: "NewField2", embeddedUnexported: embeddedUnexported{Field3: "NewField3"}},
			expected: &struct {
				Field1 string
				Field2 string
				embeddedUnexported
			}{Field1: "NewField1", Field2: "NewField2", embeddedUnexported: embeddedUnexported{Field3: "Field3"}},
		},
		{
			name: "EmbeddedStructExported",
			oldStruct: &struct {
				Field1 string
				Field2 string
				EmbeddedExported
			}{Field1: "Field1", Field2: "Field2", EmbeddedExported: EmbeddedExported{Field3: "Field3"}},
			newStruct: &struct {
				Field1 string
				Field2 string
				EmbeddedExported
			}{Field1: "NewField1", Field2: "NewField2", EmbeddedExported: EmbeddedExported{Field3: "NewField3"}},
			expected: &struct {
				Field1 string
				Field2 string
				EmbeddedExported
			}{Field1: "NewField1", Field2: "NewField2", EmbeddedExported: EmbeddedExported{Field3: "NewField3"}},
		},
		{
			name: "EmbeddedStructExportedWithPointer",
			oldStruct: &struct {
				Field1 string
				Field2 string
				*EmbeddedExported
			}{Field1: "Field1", Field2: "Field2", EmbeddedExported: &EmbeddedExported{Field3: "Field3"}},
			newStruct: &struct {
				Field1 string
				Field2 string
				*EmbeddedExported
			}{Field1: "NewField1", Field2: "NewField2", EmbeddedExported: &EmbeddedExported{Field3: "NewField3"}},
			expected: &struct {
				Field1 string
				Field2 string
				*EmbeddedExported
			}{Field1: "NewField1", Field2: "NewField2", EmbeddedExported: &EmbeddedExported{Field3: "NewField3"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			MergeStruct(test.oldStruct, test.newStruct)
			if !reflect.DeepEqual(test.oldStruct, test.expected) {
				t.Fatalf("expected %v, got %v", test.expected, test.oldStruct)
			}
		})
	}
}
