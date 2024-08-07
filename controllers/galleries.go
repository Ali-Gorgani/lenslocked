package controllers

import (
	"fmt"
	"path/filepath"

	"net/http"
	"net/url"
	"strconv"

	"github.com/Ali-Gorgani/lenslocked/context"
	"github.com/Ali-Gorgani/lenslocked/errors"
	"github.com/Ali-Gorgani/lenslocked/models"
	"github.com/go-chi/chi/v5"
)

type Galleries struct {
	Template struct {
		New   Template
		Edit  Template
		Show  Template
		Index Template
	}
	GalleryService *models.GalleryService
}

func (g Galleries) New(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title string
	}
	data.Title = r.FormValue("title")
	g.Template.New.Execute(w, r, data)
}

func (g Galleries) Create(w http.ResponseWriter, r *http.Request) {
	var data struct {
		Title  string
		UserID int
	}
	data.Title = r.FormValue("title")
	data.UserID = context.User(r.Context()).ID
	gallery, err := g.GalleryService.Create(data.UserID, data.Title)
	if err != nil {
		if errors.Is(err, models.ErrTitleTaken) {
			err = errors.Public(err, "This title is already taken")
		}
		g.Template.New.Execute(w, r, data, err)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) Edit(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       gallery.ID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Template.Edit.Execute(w, r, data)
}

func (g Galleries) Update(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	gallery.Title = r.FormValue("title")
	err = g.GalleryService.Update(gallery)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/galleries/%d", gallery.ID), http.StatusFound)
}

func (g Galleries) Show(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r)
	if err != nil {
		return
	}
	type Image struct {
		GalleryID       int
		Filename        string
		FilenameEscaped string
	}
	var data struct {
		ID     int
		Title  string
		UserID int
		Images []Image
	}
	data.ID = gallery.ID
	data.Title = gallery.Title
	data.UserID = gallery.UserID

	images, err := g.GalleryService.Images(gallery.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, image := range images {
		data.Images = append(data.Images, Image{
			GalleryID:       image.GalleryID,
			Filename:        image.Filename,
			FilenameEscaped: url.PathEscape(image.Filename),
		})
	}
	g.Template.Show.Execute(w, r, data)
}

func (g Galleries) Index(w http.ResponseWriter, r *http.Request) {
	type Gallery struct {
		ID    int
		Title string
	}
	var data struct {
		Galleries []Gallery
	}
	user := context.User(r.Context())
	galleries, err := g.GalleryService.ByUserID(user.ID)
	if err != nil {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	for _, gallery := range galleries {
		data.Galleries = append(data.Galleries, Gallery{
			ID:    gallery.ID,
			Title: gallery.Title,
		})
	}
	g.Template.Index.Execute(w, r, data)
}

func (g Galleries) Delete(w http.ResponseWriter, r *http.Request) {
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

func (g Galleries) Image(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	galleryID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return
	}
	image, err := g.GalleryService.Image(galleryID, filename)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			http.Error(w, "Image not found", http.StatusNotFound)
			return
		}
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, image.Path)
}

func (g Galleries) UploadImage(w http.ResponseWriter, r *http.Request) {
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	fileHeaders := r.MultipartForm.File["images"]
	for _, fileHeader := range fileHeaders {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
		defer file.Close()

		err = g.GalleryService.CreateImage(gallery.ID, fileHeader.Filename, file)
		if err != nil {
			var fileError models.FileError
			if errors.As(err, &fileError) {
				msg := fmt.Sprintf("%v has an invalid type or extention. Only JPG, JPEG, PNG and GIF files can be uploaded.", fileHeader.Filename)
				http.Error(w, msg, http.StatusBadRequest)
				return
			}
			fmt.Println(err)
			http.Error(w, "Something went wrong", http.StatusInternalServerError)
			return
		}
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) DeleteImage(w http.ResponseWriter, r *http.Request) {
	filename := g.filename(w, r)
	gallery, err := g.galleryByID(w, r, userMustOwnGallery)
	if err != nil {
		return
	}
	err = g.GalleryService.DeleteImage(gallery.ID, filename)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}
	editPath := fmt.Sprintf("/galleries/%d/edit", gallery.ID)
	http.Redirect(w, r, editPath, http.StatusFound)
}

func (g Galleries) filename(w http.ResponseWriter, r *http.Request) string {
	filename := chi.URLParam(r, "filename")
	filename = filepath.Base(filename)
	return filename
}

type galleryOpt func(w http.ResponseWriter, r *http.Request, gallery *models.Gallery) error

func (g Galleries) galleryByID(w http.ResponseWriter, r *http.Request, opts ...galleryOpt) (*models.Gallery, error) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid gallery ID", http.StatusBadRequest)
		return nil, err
	}
	gallery, err := g.GalleryService.ByID(id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
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
	if user.ID != gallery.UserID {
		http.Error(w, "You are not the owner of this gallery", http.StatusForbidden)
		return fmt.Errorf("user %d does not have permission to modify gallery %d", user.ID, gallery.ID)
	}
	return nil
}
