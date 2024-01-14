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

package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/praromvik/praromvik/api"
	"github.com/praromvik/praromvik/pkg/auth"
)

type User struct {
	*api.User
}

func (u *User) SignUp(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		var formData User
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		fmt.Println(formData.User)
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (u *User) SignIn(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		var formData User
		if err := json.NewDecoder(r.Body).Decode(&formData); err != nil {
			http.Error(w, "Error parsing JSON", http.StatusBadRequest)
			return
		}
		fmt.Println(formData.User)
		if formData.Valid() {
			if err := auth.GenerateJWTAndSetCookie(w, formData.User); err != nil {
				http.Error(w, fmt.Sprintf("Failed to set JWT token in cookie. %s", err), http.StatusBadRequest)
				return
			}
		}
		w.WriteHeader(http.StatusOK)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (u *User) Valid() bool {
	// TODO
	return true
}

func (u *User) SignOut(writer http.ResponseWriter, request *http.Request) {
	auth.UnsetJWTInCookie(writer, request)
	// http.Redirect(writer, request, "/signin", http.StatusSeeOther)
}
