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
	"encoding/json"
	"net/http"

	"github.com/praromvik/praromvik/models"
	"github.com/praromvik/praromvik/models/user"
	"github.com/praromvik/praromvik/pkg/auth"
	perror "github.com/praromvik/praromvik/pkg/error"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type User struct {
	*user.User
}

func (u User) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&u.User); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
			return
		}
		errCode, err := u.User.ValidateForm()
		if err != nil {
			perror.HandleError(w, errCode, "", err)
			return
		}

		u.UUID = uuid.NewString()
		if err := u.HashPassword(); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "Failed to hash password", err)
		}
		if u.Email == models.AdminEmail {
			u.Role = string(models.Admin)
		} else {
			u.Role = string(models.Student)
		}

		if err := u.User.AddUserDataToDB(); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "failed to add form data into database", err)
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (u User) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&u.User); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
			return
		}
		valid, err := u.User.VerifyLoginData()
		if err != nil && status.Code(err) != codes.NotFound {
			perror.HandleError(w, http.StatusUnauthorized, "failed to login", err)
			return
		}
		if !valid {
			perror.HandleError(w, http.StatusUnauthorized, "invalid username or password", nil)
		}
		if err := auth.StoreAuthenticated(w, r, u.User, true); err != nil {
			perror.HandleError(w, http.StatusInternalServerError, "failed to store session token", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (u User) SignOut(w http.ResponseWriter, r *http.Request) {
	if err := auth.StoreAuthenticated(w, r, u.User, false); err != nil {
		perror.HandleError(w, http.StatusInternalServerError, "failed to store session token", err)
		return
	}
}

func (u User) ProvideRoleToUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := json.NewDecoder(r.Body).Decode(&u.User); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "Error on parsing JSON", err)
			return
		}

		if err := u.User.FetchAndSetUUIDFromDB(); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "", err)
			return
		}

		if err := u.UpdateUserDataToDB(); err != nil {
			perror.HandleError(w, http.StatusBadRequest, "Error on Update User", err)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusBadRequest)
	}
}
