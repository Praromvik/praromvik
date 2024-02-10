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

package middleware

import (
	"net/http"

	"github.com/praromvik/praromvik/models"
	"github.com/praromvik/praromvik/pkg/auth"
	"github.com/praromvik/praromvik/pkg/error"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func AddMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		valid, err := auth.SessionValid(r)
		if err != nil {
			error.HandleError(w, http.StatusUnauthorized, "Failed to validate session: "+err.Error(), err)
			return
		}
		if !valid {
			error.HandleError(w, http.StatusUnauthorized, "Invalid session token", err)
			return
		}
		authenticated, err := auth.IsAuthenticated(r)
		if err != nil {
			error.HandleError(w, http.StatusUnauthorized, "Failed to check authentication status: "+err.Error(), err)
			return
		}

		if !authenticated {
			error.HandleError(w, http.StatusUnauthorized, "Unauthenticated session token", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		roleType, err := auth.GetSessionRole(request)
		if err != nil {
			error.HandleError(writer, http.StatusUnauthorized, "Failed to retrieve role from session", err)
			return
		}
		if roleType != models.Admin {
			error.HandleError(writer, http.StatusUnauthorized, "Insufficient privileges. Admin or Moderator access required", err)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func AdminOrModeratorAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		roleType, err := auth.GetSessionRole(request)
		if err != nil {
			error.HandleError(writer, http.StatusUnauthorized, "Failed to retrieve role from session", err)
			return
		}
		if !(roleType == models.Admin || roleType == models.Moderator) {
			error.HandleError(writer, http.StatusUnauthorized, "Insufficient privileges. Admin or Moderator access required", err)
			return
		}
		next.ServeHTTP(writer, request)
	})
}

func ModeratorAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		roleType, err := auth.GetSessionRole(request)
		if err != nil {
			error.HandleError(writer, http.StatusUnauthorized, "Failed to retrieve role from session", err)
			return
		}
		if roleType != models.Moderator {
			error.HandleError(writer, http.StatusUnauthorized, "Insufficient privileges. Admin or Moderator access required", err)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
