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
	"errors"
	"fmt"
	"net/http"

	"github.com/praromvik/praromvik/models/db"
	"github.com/praromvik/praromvik/models/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userMongoNamespace = db.Namespace{Database: "praromvik", Collection: "users"}

func (u *User) AddUserDataToDB() error {
	if err := u.AddUserDataToMongo(); err != nil {
		return err
	}
	if err := u.AddUserAuthDataToFirestore(); err != nil {
		return err
	}
	return nil
}

func (u *User) UpdateUserDataToDB() error {
	if err := u.UpdateUserDataToMongo(); err != nil {
		return err
	}
	if err := u.UpdateUserAuthDataToFirestore(); err != nil {
		return err
	}
	return nil
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

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(u.Password)); err != nil {
		return false, nil
	}
	u.UUID, u.Role = user.UUID, user.Role
	return true, nil
}

func (u *User) ValidateForm() (int, error) {
	keyVal := map[string]string{"userName": u.UserName, "email": u.Email, "phone": u.Phone}
	for key, val := range keyVal {
		if err := checkFieldAvailability(key, val); err != nil {
			return http.StatusBadRequest, err
		}
	}
	return http.StatusOK, nil
}

func (u *User) UpdateUserDataToMongo() error {
	user := User{}
	filter := bson.D{{Key: utils.UUID, Value: u.UUID}}
	mongoDB := db.Mongo{Namespaces: []db.Namespace{userMongoNamespace}}
	result, err := mongoDB.GetDocument(filter)
	if err != nil {
		return err
	}
	if err := result.Decode(&user); err != nil {
		return err
	}
	db.MergeStruct(&user, *u)
	_, err = mongoDB.UpdateDocument(filter, user)
	return err
}

func (u *User) AddUserDataToMongo() error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{userMongoNamespace}}
	_, err := mongoDB.AddDocument(u)
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
	err = db.Firestore{}.UpdateDocument("users", u.UserName, authData)
	return err
}

func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) FetchAndSetUUIDFromDB() error {
	var user User
	dSnap, err := db.Firestore{}.GetDocument("users", u.UserName)
	if err != nil {
		return err
	}
	if err := dSnap.DataTo(&user); err != nil {
		return err
	}
	u.UUID = user.UUID
	return nil
}

func (u *User) GetFromMongo() error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{userMongoNamespace}}
	result, err := mongoDB.GetDocument(bson.D{{Key: "userName", Value: u.UserName}})
	if err != nil {
		return err
	}
	if err := result.Decode(&u); err != nil {
		return err
	}
	return nil
}

func checkFieldAvailability(field string, value string) error {
	mongoDB := db.Mongo{Namespaces: []db.Namespace{userMongoNamespace}}
	result, err := mongoDB.GetDocument(bson.D{{Key: field, Value: value}})
	if err != nil {
		return err
	}
	if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
		return fmt.Errorf("this %s is already in use", field)
	}
	return nil
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
