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

package user

import (
	"cloud.google.com/go/firestore"
	"context"
)

func AddFormData(client *firestore.Client, formData *User) error {
	_, err := client.Collection("user").Doc(formData.UserName).Set(context.Background(), formData)
	return err
}

func VerifyLoginData(client *firestore.Client, formData *User) (bool, error) {
	dSnap, err := client.Collection("Users").Doc(formData.UserName).Get(context.Background())
	if err != nil {
		return false, err
	}

	var u User
	if err := dSnap.DataTo(&u); err != nil {
		return false, err
	}

	return u.Password == formData.Password, nil
}
