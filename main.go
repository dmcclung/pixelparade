package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func executeTemplate(w http.ResponseWriter, tPath string, tData TemplateData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	t, err := template.ParseFiles(tPath)
	if err != nil {
		log.Printf("error parsing template %v, %v\n", tPath, err)
		http.Error(w, "error parsing template", http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, tData)
	if err != nil {
		log.Printf("error executing template %v, %v\n", tPath, err)
		http.Error(w, "error executing template", http.StatusInternalServerError)
	}
}

type TemplateData struct {
	Id string
}

func galleryHandler(w http.ResponseWriter, r *http.Request) {
	galleryId := chi.URLParam(r, "id")
	executeTemplate(w, filepath.Join("templates", "gallery.gohtml"), TemplateData{
		Id: galleryId,
	})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "home.gohtml"), TemplateData{})
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "contact.gohtml"), TemplateData{})
}

func faqHandler(w http.ResponseWriter, r *http.Request) {
	executeTemplate(w, filepath.Join("templates", "faq.gohtml"), TemplateData{})
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
