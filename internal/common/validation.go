package common

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrEmptyPath indicates an empty path was provided
	ErrEmptyPath = errors.New("path cannot be empty")
	// ErrPathNotExists indicates the path does not exist
	ErrPathNotExists = errors.New("path does not exist")
	// ErrInvalidPath indicates the path is invalid
	ErrInvalidPath = errors.New("invalid path")
	// ErrNotAFile indicates the path is not a file
	ErrNotAFile = errors.New("path is not a file")
	// ErrNotADirectory indicates the path is not a directory
	ErrNotADirectory = errors.New("path is not a directory")
	// ErrEmptyString indicates an empty string was provided
	ErrEmptyString = errors.New("string cannot be empty")
)

// ValidateNonEmptyPath checks if a path is not empty
func ValidateNonEmptyPath(path string) error {
	if strings.TrimSpace(path) == "" {
		return ErrEmptyPath
	}
	return nil
}

// ValidateFileExists checks if a file exists and is not a directory
func ValidateFileExists(path string) error {
	if err := ValidateNonEmptyPath(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrPathNotExists
		}
		return fmt.Errorf("stat file: %w", err)
	}

	if info.IsDir() {
		return ErrNotAFile
	}

	return nil
}

// ValidateDirectoryExists checks if a directory exists
func ValidateDirectoryExists(path string) error {
	if err := ValidateNonEmptyPath(path); err != nil {
		return err
	}

	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ErrPathNotExists
		}
		return fmt.Errorf("stat directory: %w", err)
	}

	if !info.IsDir() {
		return ErrNotADirectory
	}

	return nil
}

// ValidatePathExists checks if a path exists (file or directory)
func ValidatePathExists(path string) error {
	if err := ValidateNonEmptyPath(path); err != nil {
		return err
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return ErrPathNotExists
		}
		return fmt.Errorf("stat path: %w", err)
	}

	return nil
}

// CleanAndValidatePath cleans and validates a path
func CleanAndValidatePath(path string) (string, error) {
	if err := ValidateNonEmptyPath(path); err != nil {
		return "", err
	}

	cleanedPath := filepath.Clean(path)
	return cleanedPath, nil
}

// ValidateNonEmptyString checks if a string is not empty
func ValidateNonEmptyString(str, fieldName string) error {
	if strings.TrimSpace(str) == "" {
		return fmt.Errorf("%s: %w", fieldName, ErrEmptyString)
	}
	return nil
}
