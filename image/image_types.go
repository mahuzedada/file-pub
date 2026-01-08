package image

import "time"

// ImageMetadata represents metadata for an uploaded image
type ImageMetadata struct {
	ID           string    `json:"id" db:"id"`
	Filename     string    `json:"filename" db:"filename"`
	OriginalName string    `json:"original_name" db:"original_name"`
	S3Key        string    `json:"s3_key" db:"s3_key"`
	S3URL        string    `json:"s3_url" db:"s3_url"`
	ContentType  string    `json:"content_type" db:"content_type"`
	Size         int64     `json:"size" db:"size"`
	UploadedAt   time.Time `json:"uploaded_at" db:"uploaded_at"`
}
