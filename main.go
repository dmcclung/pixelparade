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

func executeTemplate(w http.ResponseWriter, tPath string, tData interface{}) {
  t, err := views.Parse(tPath)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
	t.Execute(w, tData)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

  tmplt, err := views.Parse(filepath.Join("templates", "home.gohtml"))
  if err != nil {
    panic(err)
  }

  r.Get("/", controllers.Static(tmplt))

  tmplt, err = views.Parse(filepath.Join("templates", "contact.gohtml"))
  if err != nil {
    panic(err)
  }

  r.Get("/contact", controllers.Static(tmplt))

  tmplt, err = views.Parse(filepath.Join("templates", "faq.gohtml"))
  if err != nil {
    panic(err)
  }

  r.Get("/faq", controllers.Static(tmplt))

  tmplt, err = views.Parse(filepath.Join("templates", "gallery.gohtml"))
  if err != nil {
    panic(err)
  }

	r.Get("/gallery/{id}", controllers.GetGalleryById(tmplt))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})

	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
