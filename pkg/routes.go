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

package pkg

import (
	"net/http"

	"github.com/praromvik/praromvik/handlers"
	"github.com/praromvik/praromvik/pkg/middileware"

	"github.com/go-chi/chi/v5"
)

func LoadRoutes() *chi.Mux {
	router := chi.NewRouter()
	middileware.AddMiddlewares(router)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello\n"))
		w.WriteHeader(http.StatusOK)
	})

	router.Route("/courses", loadOrderRoutes)

	return router
}

func loadOrderRoutes(router chi.Router) {
	courseHandler := &handlers.Course{}

	router.Post("/", courseHandler.Create)
	router.Get("/", courseHandler.List)
	router.Get("/{id}", courseHandler.GetByID)
	router.Put("/{id}", courseHandler.UpdateByID)
	router.Delete("/{id}", courseHandler.DeleteByID)
}
