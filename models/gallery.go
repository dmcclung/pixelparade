package models

import (
	"database/sql"
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
