/*
MIT License

# Copyright (c) 2024 Praromvik

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

package routers

import (
	"github.com/praromvik/praromvik/handlers/course"
	"github.com/praromvik/praromvik/handlers/user"
	"github.com/praromvik/praromvik/pkg/auth"
	middleware "github.com/praromvik/praromvik/pkg/middileware"

	"cloud.google.com/go/firestore"
	"github.com/go-chi/chi/v5"
)

func LoadRoutes(fClient *firestore.Client) *chi.Mux {
	router := chi.NewRouter()
	middleware.AddMiddlewares(router)

	router.Group(func(r chi.Router) {
		userHandler := &user.User{FClient: fClient}
		r.HandleFunc("/signup", userHandler.SignUp)
		r.HandleFunc("/signin", userHandler.SignIn)
		r.HandleFunc("/signout", userHandler.SignOut)
	})
	router.Route("/courses", loadCourseRoutes)
	return router
}

func loadCourseRoutes(r chi.Router) {
	r.Use(auth.VerifyJWT)
	courseHandler := &course.Course{}
	r.Post("/", courseHandler.Create)
	r.Get("/", courseHandler.List)
	r.Get("/{id}", courseHandler.GetByID)
	r.Put("/{id}", courseHandler.UpdateByID)
	r.Delete("/{id}", courseHandler.DeleteByID)
}
