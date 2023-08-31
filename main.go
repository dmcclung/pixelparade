package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type GalleryData struct {
  Id string
}

func executeTemplate(w http.ResponseWriter, tPath string, tData interface{}) {
  t, err := views.Parse(tPath)
  if err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
	t.Execute(w, tData)
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	galleryId := chi.URLParam(r, "id")
  executeTemplate(w, filepath.Join("templates", "gallery.gohtml"), GalleryData{
    Id: galleryId,
  })
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
  executeTemplate(w, filepath.Join("templates", "home.gohtml"), nil)
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
  executeTemplate(w, filepath.Join("templates", "contact.gohtml"), nil)
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"), nil)
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/", homeHandler)
	r.Get("/faq", faqHandler)
	r.Get("/contact", contactHandler)
	r.Get("/gallery/{id}", galleryHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
	})
	fmt.Println("Starting the server on :3000...")
	http.ListenAndServe(":3000", r)
}
