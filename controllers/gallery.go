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
	err := g.Templates.Get.Execute(w, r, struct{ Id string }{
		Id: galleryId,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
