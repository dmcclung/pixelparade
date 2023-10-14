package controllers

import (
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

func (g Gallery) galleryByID(w http.ResponseWriter, r *http.Request) (*models.Gallery, error) {
	galleryID := chi.URLParam(r, "id")
	gallery, err := g.GalleryService.Get(galleryID)
	if err != nil {
		if err == models.ErrNoGalleryFound {
			http.Error(w, "Gallery not found", http.StatusNotFound)
			return nil, err
		}
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return nil, err
	}
	return gallery, nil
}

func (g Gallery) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to update this gallery", http.StatusUnauthorized)
		return
	}

	title := r.FormValue("title")
	gallery.Title = title

	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusNotFound)
		return
	}

	viewPath := fmt.Sprintf("/galleries/%s", gallery.ID)
	http.Redirect(w, r, viewPath, http.StatusFound)
}

func (g Gallery) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
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
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to edit this gallery", http.StatusUnauthorized)
		return
	}
	err = g.GalleryService.Delete(gallery.ID)
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

	type Gallery struct {
		ID    string
		Title string
	}

	var data struct {
		Galleries []Gallery
	}

	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}

	err = g.Templates.Index.Execute(w, r, data)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (g Gallery) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}

	var data struct {
		Title  string
		Images []string
	}
	data.Title = gallery.Title
	for i := 0; i < 10; i++ {
		data.Images = append(data.Images, fmt.Sprintf("https://placekitten.com/200/%d", 300+i))
	}

	err = g.Templates.Show.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
