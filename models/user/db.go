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
	"github.com/praromvik/praromvik/models/db"
	"go.mongodb.org/mongo-driver/bson"
)

func (u *User) UpdateUserDataToMongo() error {
	user := User{}
	dbAndCollList := []string{"praromvik", "users"}
	filter := bson.D{{Key: "userName", Value: u.UserName}}
	result, err := db.Mongo{}.GetDocument(dbAndCollList, filter)
	if err != nil {
		return err
	}
	if err := result.Decode(&user); err != nil {
		return err
	}
	db.MergeStruct(&user, *u)
	err = db.Mongo{}.UpdateDocument(dbAndCollList, filter, user)
	return err
}

func (u *User) AddUserDataToMongo() error {
	err := db.Mongo{}.AddDocument([]string{"praromvik", "users"}, u)
	return err
}

func (u *User) AddUserAuthDataToFirestore() error {
	authData := getAuthData(*u)
	err := db.Firestore{}.AddDocument("users", u.UserName, authData)
	return err
}

func (u *User) UpdateUserAuthDataToFirestore() error {
	var user User
	dSnap, err := db.Firestore{}.GetDocument("users", u.UserName)
	if err != nil {
		return err
	}
	if err := dSnap.DataTo(&user); err != nil {
		return err
	}
	db.MergeStruct(&user, *u)
	authData := getAuthData(user)
	if err != nil {
		return err
	}
	err = db.Firestore{}.UpdateDocument("users", u.UserName, authData)
	return err
}

func (u *User) VerifyLoginData() (bool, error) {
	var user User
	dSnap, err := db.Firestore{}.GetDocument("users", u.UserName)
	if err != nil {
		return false, err
	}
	if err := dSnap.DataTo(&user); err != nil {
		return false, err
	}

	if u.Password == user.Password {
		u.UUID, u.Role = user.UUID, user.Role
		return true, err
	}
	return false, nil
}

func getAuthData(user User) map[string]interface{} {
	var authData = map[string]interface{}{
		"userName": user.UserName,
		"uuid":     user.UUID,
		"password": user.Password,
		"role":     user.Role,
	}
	return authData
}
