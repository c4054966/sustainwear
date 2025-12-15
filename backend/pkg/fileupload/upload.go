package fileupload

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sustainwear/internal/config"
)

// GENERATES UNIQUE FILENAME WITH ORIGINAL EXTENSION
func GenerateUniqueFilename(originalFilename string) (string, error) {
	ext := strings.ToLower(filepath.Ext(originalFilename))

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", fmt.Errorf("failed to generate unique filename: %w", err)
	}

	uniqueID := hex.EncodeToString(randomBytes)
	filename := fmt.Sprintf("%s%s", uniqueID, ext)
	return filename, nil
}

// VALIDATES FILE EXTENSION AGAINST ALLOWED TYPES
func ValidateFileType(filename string, allowedTypes []string) error {
	ext := strings.ToLower(filepath.Ext(filename))

	for _, allowed := range allowedTypes {
		if ext == strings.ToLower(allowed) {
			return nil
		}
	}

	return fmt.Errorf("file type %s not allowed. Allowed types: %v", ext, allowedTypes)
}

// VALIDATES FILE SIZE AGAINST MAXIMUM LIMIT
func ValidateFileSize(fileSize int64, maxSizeMB int) error {
	maxSizeBytes := int64(maxSizeMB * 1024 * 1024)
	fileSizeMB := float64(fileSize) / (1024 * 1024)

	if fileSize > maxSizeBytes {
		return fmt.Errorf("file size %.2f MB exceeds maximum allowed size of %d MB", fileSizeMB, maxSizeMB)
	}

	return nil
}

// SAVES UPLOADED FILE
func SaveFile(file multipart.File, header *multipart.FileHeader, uploadDir string) (string, error) {
	uniqueFilename, err := GenerateUniqueFilename(header.Filename)
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	destPath := filepath.Join(uploadDir, uniqueFilename)
	destFile, err := os.Create(destPath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, file)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return uniqueFilename, nil
}

// VALIDATES AND SAVES FILE WITH ALL CHECKS
func ValidateAndSaveFile(file multipart.File, header *multipart.FileHeader, cfg *config.Config) (string, error) {
	if err := ValidateFileType(header.Filename, cfg.FileUpload.AllowedFileTypes); err != nil {
		return "", err
	}

	if err := ValidateFileSize(header.Size, cfg.FileUpload.MaxUploadSizeMB); err != nil {
		return "", err
	}

	filename, err := SaveFile(file, header, cfg.FileUpload.UploadDir)
	if err != nil {
		return "", err
	}

	return filename, nil
}
