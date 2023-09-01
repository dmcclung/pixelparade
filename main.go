package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/dmcclung/pixelparade/controllers"
	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

  tmplt := views.Must(views.Parse(filepath.Join("templates", "home.gohtml")))
  r.Get("/", controllers.Static(tmplt))

  tmplt = views.Must(views.Parse(filepath.Join("templates", "contact.gohtml")))
  r.Get("/contact", controllers.Static(tmplt))

  tmplt = views.Must(views.Parse(filepath.Join("templates", "faq.gohtml")))
  r.Get("/faq", controllers.Static(tmplt))

  tmplt = views.Must(views.Parse(filepath.Join("templates", "gallery.gohtml")))
	r.Get("/gallery/{id}", controllers.GetGalleryById(tmplt))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
