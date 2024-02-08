package user

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"github.com/praromvik/praromvik/models/client"
)

func (u *User) AddFormDataToMongo() error {
	coll := client.Mongo.Database("praromvik").Collection("users")
	result, err := coll.InsertOne(context.TODO(), u)

	fmt.Printf("Inserted document with _id: %v\n", result.InsertedID)
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
