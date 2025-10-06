package utils

import (
	"fmt"
	"mime"
	"path/filepath"
	"strings"

	"github.com/tdmdh/fit-up-server/internal/message/types"
)

type FileValidationConfig struct {
	MaxFileSize       int64
	AllowedMimeTypes  []string 
	AllowedExtensions []string
}


func DefaultFileValidationConfig() *FileValidationConfig {
	return &FileValidationConfig{
		MaxFileSize: 10 * 1024 * 1024, 
		AllowedMimeTypes: []string{
			"image/jpeg", "image/png", "image/webp",
			"application/pdf",
			"application/vnd.openxmlformats-officedocument.wordprocessingml.document", // .docx
			"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",       // .xlsx
		},
		AllowedExtensions: []string{
			".jpg", ".jpeg", ".png", ".webp",
			".pdf", ".docx", ".xlsx",
		},
	}
}

func ValidateFile(filename string, fileSize int64, mimeType string, config *FileValidationConfig) error {
	if config == nil {
		config = DefaultFileValidationConfig()
	}

	if fileSize > config.MaxFileSize {
		return fmt.Errorf("file size %d bytes exceeds maximum allowed size %d bytes", fileSize, config.MaxFileSize)
	}

	ext := strings.ToLower(filepath.Ext(filename))
	if len(config.AllowedExtensions) > 0 {
		allowed := false
		for _, allowedExt := range config.AllowedExtensions {
			if ext == allowedExt {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("file extension %s is not allowed", ext)
		}
	}

	if mimeType != "" && len(config.AllowedMimeTypes) > 0 {
		allowed := false
		for _, allowedType := range config.AllowedMimeTypes {
			if mimeType == allowedType {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("MIME type %s is not allowed", mimeType)
		}
	}

	return nil
}

func DetectFileType(filename string, mimeType string) types.AttachmentType {
	ext := strings.ToLower(filepath.Ext(filename))

	if mimeType != "" {
		switch {
		case strings.HasPrefix(mimeType, "image/"):
			return types.AttachmentTypeImage
		case mimeType == "application/pdf" ||
			strings.Contains(mimeType, "document") ||
			strings.Contains(mimeType, "sheet"):
			return types.AttachmentTypeDocument
		}
	}

	switch ext {
	case ".jpg", ".jpeg", ".png", ".webp":
		return types.AttachmentTypeImage
	case ".pdf", ".docx", ".xlsx":
		return types.AttachmentTypeDocument
	default:
		return types.AttachmentTypeDocument
	}
}

func GetMimeTypeFromExtension(filename string) string {
	ext := filepath.Ext(filename)
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		switch strings.ToLower(ext) {
		case ".webp":
			return "image/webp"
		default:
			return "application/octet-stream"
		}
	}
	return mimeType
}

func FormatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
