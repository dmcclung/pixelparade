package controllers

import (
	"fmt"
	"net/http"
	"path/filepath"

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

type galleryOption func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Gallery) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOption) (*models.Gallery, error) {
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

	for _, opt := range opts {
		err = opt(w, r, gallery)
		if err != nil {
			return nil, err
		}
	}

	return gallery, nil
}

func userMustOwnGallery(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error {
	user := context.User(r.Context())
	if gallery.UserID != user.ID {
		http.Error(w, "You are not authorized to update this gallery", http.StatusUnauthorized)
		return fmt.Errorf("User does not have access to this gallery")
	}
	return nil
}

func (g Gallery) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
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
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
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
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
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

	type Image struct {
		GalleryID string
		Filename  string
	}

	var data struct {
		Title  string
		Images []Image
	}
	data.Title = gallery.Title

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Print(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID: gallery.ID,
			Filename:  image.Filename,
		})
	}

	err = g.Templates.Show.Execute(w, r, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (g Gallery) Image(w http.ResponseWriter, r *http.Request) {
	galleryID := chi.URLParam(r, "id")
	filename := chi.URLParam(r, "filename")

	path := filepath.Join(g.GalleryService.GalleryDir(galleryID), filename)
	http.ServeFile(w, r, path)
}
