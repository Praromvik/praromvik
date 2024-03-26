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
	"context"
	"net/http"

	"github.com/praromvik/praromvik/models"
	"github.com/praromvik/praromvik/models/course"
	"github.com/praromvik/praromvik/pkg/auth"
	perror "github.com/praromvik/praromvik/pkg/error"

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
			perror.HandleError(w, http.StatusUnauthorized, "Failed to validate session: "+err.Error(), err)
			return
		}
		if !valid {
			perror.HandleError(w, http.StatusUnauthorized, "Invalid session token", err)
			return
		}
		authenticated, err := auth.IsAuthenticated(r)
		if err != nil {
			perror.HandleError(w, http.StatusUnauthorized, "Failed to check authentication status: "+err.Error(), err)
			return
		}

		if !authenticated {
			perror.HandleError(w, http.StatusUnauthorized, "Unauthenticated session token", err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func AdminAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		serveHTTPIfRoleMatched(next, writer, request, []models.RoleType{models.Admin})
	})
}

func AdminOrModeratorAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		serveHTTPIfRoleMatched(next, writer, request, []models.RoleType{models.Moderator, models.Admin})
	})
}

func ModeratorAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		serveHTTPIfRoleMatched(next, writer, request, []models.RoleType{models.Moderator})
	})
}

func serveHTTPIfRoleMatched(next http.Handler, writer http.ResponseWriter, request *http.Request, roles []models.RoleType) {
	role, err := auth.GetSessionRole(request)
	if err != nil {
		perror.HandleError(writer, http.StatusUnauthorized, "Failed to retrieve role from session", err)
		return
	}
	var authenticated bool
	for _, val := range roles {
		if val == role {
			authenticated = true
			break
		}
	}
	if !authenticated {
		perror.HandleError(writer, http.StatusUnauthorized, "Insufficient privileges.", err)
	}
	next.ServeHTTP(writer, request)
}

func AddCourseUUIDToCtx(next http.Handler) http.Handler {
	type ctxKey string
	const keyUUID ctxKey = "uuid"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		uuid, err := course.GetCourseUUID(id)
		if err != nil {
			perror.HandleError(w, http.StatusNotFound, "", err)
			return
		}
		// Add UUID to request context
		ctx := context.WithValue(r.Context(), keyUUID, uuid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
