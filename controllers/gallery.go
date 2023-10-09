package controllers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dmcclung/pixelparade/context"
	"github.com/dmcclung/pixelparade/models"
	"github.com/go-chi/chi/v5"
)

type GalleryTemplates struct {
	New   Template
	Edit  Template
	Show  Template
	Index Template
}

type Gallery struct {
	Templates      GalleryTemplates
	GalleryService *models.GalleryService
}

func (g Gallery) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	err := g.Templates.New.Execute(w, r, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (g Gallery) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		UserID string
		Title  string
	}
	data.UserID = context.User(r.Context()).ID
	data.Title = r.FormValue("title")

	_, err := g.GalleryService.Create(data.Title, data.UserID)
	if err != nil {
		g.Templates.New.Execute(w, r, data, err)
		return
	}

	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Gallery) Update(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")
	gallery, err := g.GalleryService.Get(galleryID)
	if err != nil {
		if err == models.ErrNoGalleryFound {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to update this gallery", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")

	updatedGallery, err := g.GalleryService.Update(galleryID, title)
	if err != nil {
		if err == models.ErrNoGalleryFound {
			http.Error(w, "No gallery found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusNotFound)
		return
	}

	viewPath := fmt.Sprintf("/galleries/%s", updatedGallery.ID)
	http.Redirect(w, r, viewPath, http.StatusFound)
}

func (g Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")
	gallery, err := g.GalleryService.Get(galleryID)
	if err != nil {
		if err == models.ErrNoGalleryFound {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusUnauthorized)
		return
	}

	data := struct {
		ID    string
		Title string
	}{
		ID:    gallery.ID,
		Title: gallery.Title,
	}

	g.Templates.Edit.Execute(w, r, data)
}

func (g Gallery) Delete(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")
	gallery, err := g.GalleryService.Get(galleryID)
	if err != nil {
		if err == models.ErrNoGalleryFound {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusUnauthorized)
		return
	}
	err = g.GalleryService.Delete(galleryID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/galleries", http.StatusFound)
}

func (g Gallery) Index(w http.ResponseWriter, r *http.Request) {
	user := context.User(r.Context())
	galleries, err := g.GalleryService.GetByUser(user.ID)
	if err != nil {
		if err == models.ErrNoGalleries {
			galleries = []*models.Gallery{}
		} else {
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
		}
	}

	var data struct{
		Galleries []*models.Gallery
	}
	data.Galleries = galleries

	err = g.Templates.Index.Execute(w, r, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (g Gallery) Show(w http.ResponseWriter, r *http.Request) {
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

	err = g.Templates.Show.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
