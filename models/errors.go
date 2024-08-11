package models

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
)

var (
	ErrNotFound   = errors.New("models: resource could not be found")
	ErrEmailTaken = errors.New("models: email address is already taken")
	ErrTitleTaken = errors.New("models: gallery title is already taken")
)

type FileError struct {
	Issue string
}

func (fe FileError) Error() string {
	return fe.Issue
}

func checkContentType(r io.Reader, allowedTypes []string) ([]byte, error) {
	// Read the first 512 bytes to detect the content type
	buffer := make([]byte, 512)
	n, err := r.Read(buffer)
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read the first 512 bytes: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	for _, allowedType := range allowedTypes {
		if contentType == allowedType {
			return buffer[:n], nil
		}
	}
	return nil, FileError{
		Issue: fmt.Sprintf("invalid file type: %s", contentType),
	}
}

func checkExtension(filename string, allowedExtensions []string) error {
	if !hasExtention(filename, allowedExtensions) {
		return FileError{
			Issue: fmt.Sprintf("invalid file extension: %s", filepath.Ext(filename)),
		}
	}
	return nil
}
