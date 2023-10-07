package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Gallery struct {
	ID     string
	UserID string
	Title  string
}

type GalleryService struct {
	DB *sql.DB
}

func (gs *GalleryService) Create(title, userID string) (*Gallery, error) {
	gallery := Gallery{
		Title: title,
		UserID: userID,
	}
	err := gs.DB.QueryRow(`
		INSERT INTO galleries (userID, title) VALUES ($1, $2)
		RETURNING id;`, title, userID).Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (gs *GalleryService) Get(galleryID string) (*Gallery, error) {
	return nil, errors.ErrUnsupported
}

func (gs *GalleryService) GetByUser(userID string) ([]*Gallery, error) {
	return nil, errors.ErrUnsupported
}

func (gs *GalleryService) Update(title string) (*Gallery, error) {
	return nil, errors.ErrUnsupported
}

func (gs *GalleryService) Delete(galleryID string) error {
	return errors.ErrUnsupported
}
