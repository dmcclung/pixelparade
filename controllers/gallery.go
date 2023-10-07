package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Gallery struct {
	Templates struct {
		Get Template
	}
}

func (g Gallery) Get(w http.ResponseWriter, r *http.Request) {
	galleryId := chi.URLParam(r, "id")

	var data struct {
		Id string
	}
	data.Id = galleryId

	err := g.Templates.Get.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
