package fdb

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/praromvik/praromvik/api"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

func FireStoreConnect(ctx context.Context) (*firestore.Client, error) {
	// Use a service account
	SAPath := os.Getenv("FIRESTORE_SERVICE_ACCOUNT_JSON_KEY")
	opt := option.WithCredentialsFile(SAPath)

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return client, nil
}

func AddFormData(client *firestore.Client, formData *api.User) error {
	_, err := client.Collection("Users").Doc(formData.UserName).Set(context.Background(), formData)
	if err != nil {
		return err
	}
	return nil
}

func VerifyLoginData(client *firestore.Client, formData *api.User) (bool, error) {
	dsnap, err := client.Collection("Users").Doc(formData.UserName).Get(context.Background())
	if err != nil {
		return false, err
	}

	var u api.User
	if err := dsnap.DataTo(&u); err != nil {
		return false, err
	}

	return u.Password == formData.Password, nil
}
