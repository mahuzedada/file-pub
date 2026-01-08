package image

import "errors"

var (
	// ErrInvalidImageType indicates an invalid image file type was provided
	ErrInvalidImageType = errors.New("invalid image type, only JPEG, PNG, GIF, and WebP are allowed")
	// ErrImageNotFound indicates the requested image was not found
	ErrImageNotFound = errors.New("image not found")
	// ErrFileTooLarge indicates the uploaded file is too large
	ErrFileTooLarge = errors.New("file too large")
)
