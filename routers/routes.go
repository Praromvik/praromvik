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
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/praromvik/praromvik/handlers/course"
	"github.com/praromvik/praromvik/handlers/user"
	middleware "github.com/praromvik/praromvik/pkg/middileware"
)

func LoadRoutes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// Apply global middleware
	middleware.AddMiddlewares(router)

	router.Group(func(r chi.Router) {
		loadUserAuthRoutes(r)
	})
	router.With(middleware.AdminAccess).Post("/role", user.User{}.ProvideRoleToUser)

	router.Route("/course", loadCourseRoutes)
	return router
}

func loadUserAuthRoutes(r chi.Router) {
	userHandler := &user.User{}
	r.Put("/signup", userHandler.SignUp)
	r.Put("/signin", userHandler.SignIn)
	r.Delete("/signout", userHandler.SignOut)
	r.With(middleware.SecurityMiddleware).Get("/user/{userName}", userHandler.Get)
}

func loadCourseRoutes(r chi.Router) {
	r.Use(middleware.SecurityMiddleware)
	//r.Route("/lesson", loadLessonRoutes)
	handler := &course.Course{}

	r.Get("/list", handler.List)
	r.Get("/{id}", handler.Get)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AdminOrModeratorAccess)
		r.Post("/", handler.Create)
		//r.With(middleware.AddCourseIDToCtx).Put("/{name}", handler.Update)
	})

	// Require admin access and use AddCourseIDToCtx middleware
	//r.With(middleware.AdminAccess, middleware.AddCourseIDToCtx).Delete("/{name}", handler.Delete)

}

//func loadLessonRoutes(r chi.Router) {
//	handler := &course.Lesson{}
//	r.Get("/list", handler.List)
//	r.With(middleware.AddCourseIDToCtx).Get("/{name}", handler.Get)
//
//	r.Group(func(r chi.Router) {
//		r.Use(middleware.AdminOrModeratorAccess)
//		r.Post("/", handler.Create)
//		r.With(middleware.AddCourseIDToCtx).Put("/{name}", handler.Update)
//	})
//
//	// Require admin access and use AddCourseIDToCtx middleware
//	r.With(middleware.AdminAccess, middleware.AddCourseIDToCtx).Delete("/{name}", handler.Delete)
//
//}
