package models

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Image struct {
	GalleryID int
	Path      string
	Filename  string
}

type Gallery struct {
	ID     int
	UserID int
	Title  string
}

type GalleryService struct {
	DB *sql.DB

	// ImageDir represents the directory path where images for this gallery are stored
	// If not set, the GalleryService will use the default image directory "images" in the current working directory
	ImageDir string
}

func (service *GalleryService) Create(userID int, title string) (*Gallery, error) {
	gallery := Gallery{
		UserID: userID,
		Title:  title,
	}
	row := service.DB.QueryRow(`
		SELECT id FROM galleries
		WHERE title = $1 AND user_id = $2;`, gallery.Title, gallery.UserID)
	err := row.Scan(&gallery.ID)
	if err == nil {
		return nil, fmt.Errorf("create gallery: %w", ErrTitleTaken)
	}

	row = service.DB.QueryRow(`
		INSERT INTO galleries (user_id, title)
		VALUES ($1, $2) RETURNING id;`,
		gallery.UserID,
		gallery.Title,
	)
	err = row.Scan(&gallery.ID)
	if err != nil {
		return nil, fmt.Errorf("create gallery: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByID(id int) (*Gallery, error) {
	gallery := Gallery{
		ID: id,
	}
	row := service.DB.QueryRow(`
		SELECT title, user_id
		FROM galleries
		WHERE id = $1;`,
		gallery.ID,
	)
	err := row.Scan(&gallery.Title, &gallery.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("query gallery by id: %w", err)
	}
	return &gallery, nil
}

func (service *GalleryService) ByUserID(userID int) ([]Gallery, error) {
	rows, err := service.DB.Query(`
		SELECT id, title
		FROM galleries
		WHERE user_id = $1
		ORDER BY id;`, userID)
	if err != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	var galleries []Gallery
	for rows.Next() {
		gallery := Gallery{
			UserID: userID,
		}
		err = rows.Scan(&gallery.ID, &gallery.Title)
		if err != nil {
			return nil, fmt.Errorf("query galleries by user: %w", err)
		}
		galleries = append(galleries, gallery)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("query galleries by user: %w", err)
	}
	return galleries, nil
}

func (service *GalleryService) Update(gallery *Gallery) error {
	_, err := service.DB.Exec(`
		UPDATE galleries
		SET title = $2
		WHERE id = $1;`,
		gallery.ID,
		gallery.Title,
	)
	if err != nil {
		return fmt.Errorf("update gallery title: %w", err)
	}
	return err
}

func (service *GalleryService) Delete(id int) error {
	_, err := service.DB.Exec(`
		DELETE FROM galleries
		WHERE id = $1;`,
		id,
	)
	if err != nil {
		return fmt.Errorf("delete gallery: %w", err)
	}

	galleryDir := service.galleryDir(id)
	err = os.RemoveAll(galleryDir)
	if err != nil {
		return fmt.Errorf("delete gallery directory: %w", err)
	}

	return err
}

func (service *GalleryService) Images(galleryID int) ([]Image, error) {
	globPattern := filepath.Join(service.galleryDir(galleryID), "*")
	allFiles, err := filepath.Glob(globPattern)
	if err != nil {
		return nil, fmt.Errorf("retrieving gallery images: %w", err)
	}
	var images []Image
	for _, file := range allFiles {
		if hasExtention(file, service.extensions()) {
			images = append(images, Image{
				GalleryID: galleryID,
				Path:      file,
				Filename:  filepath.Base(file),
			})
		}
	}
	return images, nil
}

func (service *GalleryService) Image(galleryID int, filename string) (Image, error) {
	imagePath := filepath.Join(service.galleryDir(galleryID), filename)
	_, err := os.Stat(imagePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return Image{}, ErrNotFound
		}
		return Image{}, fmt.Errorf("retrieving image: %w", err)
	}
	return Image{
		GalleryID: galleryID,
		Path:      imagePath,
		Filename:  filename,
	}, nil
}

func (service *GalleryService) CreateImage(galleryID int, filename string, contents io.ReadSeeker) error {
	err := checkContentType(contents, service.imageContentTypes())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}
	err = checkExtension(filename, service.extensions())
	if err != nil {
		return fmt.Errorf("creating image %v: %w", filename, err)
	}

	galleryDir := service.galleryDir(galleryID)
	err = os.MkdirAll(galleryDir, 0755)
	if err != nil {
		return fmt.Errorf("creating gallery-%d images directory: %w", galleryID, err)
	}
	imagePath := filepath.Join(galleryDir, filename)
	dst, err := os.Create(imagePath)
	if err != nil {
		return fmt.Errorf("creating image file: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, contents)
	if err != nil {
		return fmt.Errorf("copying contents to image: %w", err)
	}
	return nil
}

func (service *GalleryService) DeleteImage(galleryID int, filename string) error {
	image, err := service.Image(galleryID, filename)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	err = os.Remove(image.Path)
	if err != nil {
		return fmt.Errorf("deleting image: %w", err)
	}
	return nil
}

func (service *GalleryService) extensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".gif"}
}

func (service *GalleryService) imageContentTypes() []string {
	return []string{"image/jpeg", "image/png", "image/gif"}
}

func (service *GalleryService) galleryDir(id int) string {
	if service.ImageDir == "" {
		service.ImageDir = "images"
	}
	return filepath.Join(service.ImageDir, fmt.Sprintf("gallery-%d", id))
}

func hasExtention(file string, extention []string) bool {
	for _, ext := range extention {
		file = strings.ToLower(file)
		ext = strings.ToLower(ext)
		if filepath.Ext(file) == ext {
			return true
		}
	}
	return false
}
