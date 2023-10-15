package models

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"strings"
)

type Gallery struct {
	ID     string
	UserID string
	Title  string
}

type GalleryService struct {
	DB *sql.DB
	ImagesDir string
}

type Image struct {
	Path string
}

func hasExtension(file string, extensions []string) bool {
	for _, ext := range extensions {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}

func (gs *GalleryService) Images(id string) ([]Image, error) {
	globPath := filepath.Join(gs.galleryDir(id), "*")
	paths, err := filepath.Glob(globPath)
	if err != nil {
		return nil, fmt.Errorf("list images: %w", err)
	}

	var images []Image
	for _, path := range paths {
		if hasExtension(path, gs.extensions()) {
			images = append(images, Image{
				Path: path,
			})
		}
	}

	return images, nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".jpg", ".png", ".jpeg", ".gif"}
}

func (gs *GalleryService) galleryDir(id string) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%s", id))
}

func (gs *GalleryService) Create(title, userID string) (*Gallery, error) {
	gallery := Gallery{
		Title:  title,
		UserID: userID,
	}
	err := gs.DB.QueryRow(`
		INSERT INTO galleries (user_id, title) VALUES ($1, $2)
		RETURNING id;`, userID, title).Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) Get(galleryID string) (*Gallery, error) {
	gallery := Gallery{
		ID: galleryID,
	}
	err := gs.DB.QueryRow(`
		SELECT user_id, title FROM galleries WHERE id = $1;
	`, galleryID).Scan(&gallery.UserID, &gallery.Title)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoGalleryFound
		}
		return nil, fmt.Errorf("get gallery: %w", err)
	}

	return &gallery, nil
}

func (gs *GalleryService) GetByUser(userID string) ([]*Gallery, error) {
	rows, err := gs.DB.Query(`
		SELECT id, user_id, title FROM galleries WHERE user_id = $1;
	`, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNoGalleries
		}
		return nil, fmt.Errorf("get user galleries: %w", err)
	}

	var galleries []*Gallery
	for rows.Next() {
		var gallery Gallery
		err := rows.Scan(&gallery.ID, &gallery.UserID, &gallery.Title)
		if err != nil {
			return galleries, err
		}
		galleries = append(galleries, &gallery)
	}
	return galleries, nil
}

func (gs *GalleryService) Update(gallery *Gallery) error {
	_, err := gs.DB.Exec(`
		UPDATE galleries SET title = $2 WHERE id = $1;
	`, gallery.ID, gallery.Title)
	if err != nil {
		return fmt.Errorf("update gallery: %w", err)
	}

	return nil
}

func (gs *GalleryService) Delete(galleryID string) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries WHERE id = $1;
	`, galleryID)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	return nil
}
