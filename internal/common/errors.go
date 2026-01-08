package common

import "fmt"

// WrapFileError wraps file operation errors with context
func WrapFileError(operation, path string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s file %s: %w", operation, path, err)
}

// WrapDatabaseError wraps database operation errors with context
func WrapDatabaseError(operation string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("database %s: %w", operation, err)
}

// WrapS3Error wraps S3 operation errors with context
func WrapS3Error(operation, bucket, key string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("s3 %s bucket=%s key=%s: %w", operation, bucket, key, err)
}
