package models

import (
	"database/sql"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/dmcclung/pixelparade/errors"
)

type Gallery struct {
	ID     string
	UserID string
	Title  string
}

type GalleryService struct {
	DB        *sql.DB
	ImagesDir string
}

type Image struct {
	Path      string
	GalleryID string
	Filename  string
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
				Path:      path,
				Filename:  filepath.Base(path),
				GalleryID: id,
			})
		}
	}

	return images, nil
}

func (gs *GalleryService) extensions() []string {
	return []string{".jpg", ".png", ".jpeg", ".gif"}
}

func (gs *GalleryService) imageContentTypes() []string {
	return []string{"image/png", "image/jpeg", "image/gif"}
}

func (gs *GalleryService) galleryDir(id string) string {
	imagesDir := gs.ImagesDir
	if imagesDir == "" {
		imagesDir = "images"
	}
	return filepath.Join(imagesDir, fmt.Sprintf("gallery-%s", id))
}

func (gs *GalleryService) ImagePath(galleryID, filename string) (string, error) {
	path := filepath.Join(gs.galleryDir(galleryID), filename)

	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return "", ErrImageNotFound
		}
		return "", fmt.Errorf("image exists: %w", err)
	}

	return path, nil
}

func (gs *GalleryService) DeleteImage(galleryID, filename string) error {
	path, err := gs.ImagePath(galleryID, filename)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	err = os.Remove(path)
	if err != nil {
		return fmt.Errorf("delete image: %w", err)
	}
	return nil
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
		if errors.Is(err, sql.ErrNoRows) {
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
		if errors.Is(err, sql.ErrNoRows) {
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

func checkContentType(r io.ReadSeeker, allowedContentTypes []string) error {
	buf := make([]byte, 512)
	_, err := r.Read(buf)
	if err != nil {
		return fmt.Errorf("detecting content type: %w", err)
	}
	contentType := http.DetectContentType(buf)
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return fmt.Errorf("detecting content type: %w", err)
	}

	for _, allowedContentType := range allowedContentTypes {
		if contentType == allowedContentType {
			return nil
		}
	}

	return FileError{
		Issue: fmt.Sprintf("found %s content type, but expected %v", contentType, allowedContentTypes),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtension(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid extension: %v", filepath.Ext(filename)),
		}
	}
	return nil
}

func (gs *GalleryService) CreateImage(galleryID, filename string, file io.ReadSeeker) error {
	err := checkContentType(file, gs.imageContentTypes())
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}

	err = checkExtension(filename, gs.extensions())
	if err != nil {
		return fmt.Errorf("create image: %w", err)
	}

	galleryDir := gs.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("create gallery dir: %w", err)
	}

	dstPath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("create image: %w", err)
	}
	defer dst.Close()
	io.Copy(dst, file)

	return nil
}

func (gs *GalleryService) Delete(galleryID string) error {
	_, err := gs.DB.Exec(`
		DELETE FROM galleries WHERE id = $1;
	`, galleryID)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	err = os.RemoveAll(gs.galleryDir(galleryID))
	if err != nil {
		return fmt.Errorf("delete gallery images: %w", err)
	}

	return nil
}
