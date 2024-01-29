package user

import (
	"net/http"

	"github.com/praromvik/praromvik/api"
	fdb "github.com/praromvik/praromvik/pkg/db/firestore"
)

func (u *User) validate(w http.ResponseWriter) bool {
	var valid bool
	valid = true
	// TODO
	if !valid {
		api.HandleError(w, http.StatusConflict, "account with this email address already exists", nil)
	}
	return valid
}

func (u *User) verify() (bool, error) {
	valid, err := fdb.VerifyLoginData(u.FClient, u.User)
	if err != nil {
		return false, err
	}
	return valid, nil
}
