package image

import (
	"context"
	"database/sql"
	"fmt"

	"file-pub/internal/common"
)

// ImageRepository defines the interface for image data access
type ImageRepository interface {
	GetAllImages(ctx context.Context) ([]ImageMetadata, error)
	SaveImage(ctx context.Context, metadata ImageMetadata) error
	GetImageByID(ctx context.Context, id string) (*ImageMetadata, error)
}

// imageRepository implements ImageRepository
type imageRepository struct {
	db *sql.DB
}

// NewImageRepository creates a new ImageRepository
func NewImageRepository(db *sql.DB) ImageRepository {
	common.RequireNonNil(db, "db")

	return &imageRepository{
		db: db,
	}
}

// GetAllImages retrieves all images from the database
func (repo *imageRepository) GetAllImages(ctx context.Context) ([]ImageMetadata, error) {
	query := `
		SELECT id, filename, original_name, s3_key, s3_url, content_type, size, uploaded_at
		FROM images
		ORDER BY uploaded_at DESC
	`

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, common.WrapDatabaseError("query images", err)
	}
	defer rows.Close()

	var images []ImageMetadata
	for rows.Next() {
		var img ImageMetadata
		err := rows.Scan(
			&img.ID,
			&img.Filename,
			&img.OriginalName,
			&img.S3Key,
			&img.S3URL,
			&img.ContentType,
			&img.Size,
			&img.UploadedAt,
		)
		if err != nil {
			return nil, common.WrapDatabaseError("scan image row", err)
		}
		images = append(images, img)
	}

	if err := rows.Err(); err != nil {
		return nil, common.WrapDatabaseError("iterate image rows", err)
	}

	return images, nil
}

// SaveImage saves image metadata to the database
func (repo *imageRepository) SaveImage(ctx context.Context, metadata ImageMetadata) error {
	query := `
		INSERT INTO images (id, filename, original_name, s3_key, s3_url, content_type, size, uploaded_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := repo.db.ExecContext(
		ctx,
		query,
		metadata.ID,
		metadata.Filename,
		metadata.OriginalName,
		metadata.S3Key,
		metadata.S3URL,
		metadata.ContentType,
		metadata.Size,
		metadata.UploadedAt,
	)

	if err != nil {
		return common.WrapDatabaseError("insert image", err)
	}

	return nil
}

// GetImageByID retrieves an image by ID from the database
func (repo *imageRepository) GetImageByID(ctx context.Context, id string) (*ImageMetadata, error) {
	query := `
		SELECT id, filename, original_name, s3_key, s3_url, content_type, size, uploaded_at
		FROM images
		WHERE id = ?
	`

	var img ImageMetadata
	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&img.ID,
		&img.Filename,
		&img.OriginalName,
		&img.S3Key,
		&img.S3URL,
		&img.ContentType,
		&img.Size,
		&img.UploadedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrImageNotFound
		}
		return nil, common.WrapDatabaseError(fmt.Sprintf("query image %s", id), err)
	}

	return &img, nil
}
