package image

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"file-pub/internal/common"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

var (
	// validImageTypes defines allowed image content types
	validImageTypes = map[string]bool{
		"image/jpeg": true,
		"image/jpg":  true,
		"image/png":  true,
		"image/gif":  true,
		"image/webp": true,
	}
)

// ImageService defines the interface for image business logic
type ImageService interface {
	GetAllImages(ctx context.Context) ([]ImageMetadata, error)
	UploadImage(ctx context.Context, file io.Reader, filename, contentType string, size int64) (*ImageMetadata, error)
	ValidateImageType(contentType string) error
}

// imageService implements ImageService
type imageService struct {
	imageRepo ImageRepository
	uploader  *s3manager.Uploader
	s3Bucket  string
}

// NewImageService creates a new ImageService
func NewImageService(
	imageRepo ImageRepository,
	uploader *s3manager.Uploader,
	s3Bucket string,
) ImageService {
	common.PanicOnInvalidDependencies("ImageService", map[string]interface{}{
		"imageRepo": imageRepo,
		"uploader":  uploader,
	})

	if err := common.ValidateNonEmptyString(s3Bucket, "s3Bucket"); err != nil {
		panic(fmt.Sprintf("ImageService: %v", err))
	}

	return &imageService{
		imageRepo: imageRepo,
		uploader:  uploader,
		s3Bucket:  s3Bucket,
	}
}

// GetAllImages retrieves all images
func (service *imageService) GetAllImages(ctx context.Context) ([]ImageMetadata, error) {
	images, err := service.imageRepo.GetAllImages(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting all images: %w", err)
	}

	return images, nil
}

// UploadImage uploads an image to S3 and saves metadata to database
func (service *imageService) UploadImage(ctx context.Context, file io.Reader, filename, contentType string, size int64) (*ImageMetadata, error) {
	// Validate content type
	if err := service.ValidateImageType(contentType); err != nil {
		return nil, err
	}

	// Generate unique ID and filename
	id := uuid.New().String()
	ext := filepath.Ext(filename)
	uniqueFilename := id + ext
	s3Key := "uploads/" + uniqueFilename

	// Upload to S3
	uploadResult, err := service.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket:      aws.String(service.s3Bucket),
		Key:         aws.String(s3Key),
		Body:        file,
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, common.WrapS3Error("upload", service.s3Bucket, s3Key, err)
	}

	// Create metadata
	metadata := ImageMetadata{
		ID:           id,
		Filename:     uniqueFilename,
		OriginalName: filename,
		S3Key:        s3Key,
		S3URL:        uploadResult.Location,
		ContentType:  contentType,
		Size:         size,
		UploadedAt:   time.Now(),
	}

	// Save metadata to database
	if err := service.imageRepo.SaveImage(ctx, metadata); err != nil {
		return nil, fmt.Errorf("saving image metadata: %w", err)
	}

	return &metadata, nil
}

// ValidateImageType validates if the content type is an allowed image type
func (service *imageService) ValidateImageType(contentType string) error {
	if !validImageTypes[contentType] {
		return ErrInvalidImageType
	}
	return nil
}
