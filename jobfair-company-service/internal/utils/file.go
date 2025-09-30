package utils

import (
	"errors"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

type FileUploadConfig struct {
	AllowedExtensions []string
	MaxFileSize       int64
	UploadPath        string
}

var (
	ImageConfig = FileUploadConfig{
		AllowedExtensions: []string{".jpg", ".jpeg", ".png", ".gif", ".webp"},
		MaxFileSize:       10 * 1024 * 1024,
		UploadPath:        "/uploads/images",
	}

	VideoConfig = FileUploadConfig{
		AllowedExtensions: []string{".mp4", ".avi", ".mov", ".webm"},
		MaxFileSize:       50 * 1024 * 1024,
		UploadPath:        "/uploads/videos",
	}

	DocumentConfig = FileUploadConfig{
		AllowedExtensions: []string{".pdf", ".doc", ".docx"},
		MaxFileSize:       5 * 1024 * 1024,
		UploadPath:        "/uploads/documents",
	}
)

func ValidateFile(file *multipart.FileHeader, config FileUploadConfig) error {
	if file.Size > config.MaxFileSize {
		return fmt.Errorf("file size exceeds maximum allowed size of %dMB", config.MaxFileSize/(1024*1024))
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	valid := false
	for _, allowedExt := range config.AllowedExtensions {
		if ext == allowedExt {
			valid = true
			break
		}
	}

	if !valid {
		return errors.New("file type not allowed")
	}

	return nil
}

func GenerateFileName(originalName string, prefix string) string {
	ext := filepath.Ext(originalName)
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s_%d%s", prefix, timestamp, ext)
}

func GetFileExtension(filename string) string {
	return strings.ToLower(filepath.Ext(filename))
}

func IsImageFile(filename string) bool {
	ext := GetFileExtension(filename)
	for _, allowed := range ImageConfig.AllowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}

func IsVideoFile(filename string) bool {
	ext := GetFileExtension(filename)
	for _, allowed := range VideoConfig.AllowedExtensions {
		if ext == allowed {
			return true
		}
	}
	return false
}