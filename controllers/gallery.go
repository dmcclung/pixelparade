package controllers

import (
	"errors"
	"net/http"

	"github.com/dmcclung/pixelparade/models"
	"github.com/go-chi/chi/v5"
)

type GalleryTemplates struct {
	NewGallery   Template
	GetGalleries Template
	GetGallery   Template
}

type Gallery struct {
	Templates      GalleryTemplates
	GalleryService *models.GalleryService
}

func (g Gallery) New(w http.ResponseWriter, r *http.Request) {

}

func (g Gallery) GetGalleries(w http.ResponseWriter, r *http.Request) {

}

func (g Gallery) ProcessNewGallery(w http.ResponseWriter, r *http.Request) {

}

func (g Gallery) Get(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")

	gallery, err := g.GalleryService.Get(galleryID)
	if err != nil {
		// TODO: return home?
		if errors.Is(err, models.ErrNoGalleryFound) {
			http.Error(w, "No gallery found", http.StatusNotFound)
		}
	}

	var data struct {
		Title string
	}
	data.Title = gallery.Title

	err = g.Templates.GetGallery.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
