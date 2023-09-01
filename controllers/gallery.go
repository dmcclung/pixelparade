package controllers

import (
	"net/http"

	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
)

type Gallery struct {
	HtmlTmpl views.Template
}

type GalleryData struct {
	Id string
}

func (g Gallery) GetGalleryById(w http.ResponseWriter, r *http.Request) {
	galleryId := chi.URLParam(r, "id")
	err := g.HtmlTmpl.Execute(w, GalleryData{
		Id: galleryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
