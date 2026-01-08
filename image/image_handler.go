package image

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"file-pub/internal/common"
)

// ImageHandler handles HTTP requests for image operations
type ImageHandler struct {
	imageService ImageService
	templates    *template.Template
}

// NewImageHandler creates a new ImageHandler
func NewImageHandler(
	imageService ImageService,
	templates *template.Template,
) *ImageHandler {
	common.PanicOnInvalidDependencies("ImageHandler", map[string]interface{}{
		"imageService": imageService,
		"templates":    templates,
	})

	return &ImageHandler{
		imageService: imageService,
		templates:    templates,
	}
}

// HandleHome displays the home page with all images
func (handler *ImageHandler) HandleHome(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Fetch all images from database
	images, err := handler.imageService.GetAllImages(r.Context())
	if err != nil {
		log.Printf("Error fetching images: %v", err)
		http.Error(w, "Failed to fetch images", http.StatusInternalServerError)
		return
	}

	data := struct {
		Images []ImageMetadata
		Count  int
	}{
		Images: images,
		Count:  len(images),
	}

	if err := handler.templates.ExecuteTemplate(w, "index.html", data); err != nil {
		log.Printf("Template error: %v", err)
		http.Error(w, "Failed to render page", http.StatusInternalServerError)
	}
}

// HandleUpload handles image upload requests
func (handler *ImageHandler) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (max 32MB)
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		log.Printf("Error parsing form: %v", err)
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Error reading file: %v", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Get content type
	contentType := header.Header.Get("Content-Type")

	// Validate file type
	if err := handler.imageService.ValidateImageType(contentType); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Upload image
	_, err = handler.imageService.UploadImage(
		r.Context(),
		file,
		header.Filename,
		contentType,
		header.Size,
	)
	if err != nil {
		log.Printf("Upload error: %v", err)
		http.Error(w, fmt.Sprintf("Failed to upload file: %v", err), http.StatusInternalServerError)
		return
	}

	// Redirect back to home
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
