package controllers

import (
	"net/http"

	"github.com/dmcclung/pixelparade/views"
	"github.com/go-chi/chi/v5"
)

type GalleryData struct {
	Id string
}

func GetGalleryById(tmplt views.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		galleryId := chi.URLParam(r, "id")
		err := tmplt.Execute(w, GalleryData{
			Id: galleryId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
