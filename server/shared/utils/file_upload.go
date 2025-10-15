package utils

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	MaxFileSize       = 10 << 20 // 10 MB
	AllowedImageTypes = ".jpg,.jpeg,.png,.gif,.webp"
)

type FileUploadConfig struct {
	UploadDir      string
	MaxSize        int64
	AllowedTypes   string
	FileNamePrefix string
}

func HandleFileUpload(file multipart.File, header *multipart.FileHeader, config FileUploadConfig) (string, error) {
	if config.MaxSize == 0 {
		config.MaxSize = MaxFileSize
	}
	if config.AllowedTypes == "" {
		config.AllowedTypes = AllowedImageTypes
	}

	if header.Size > config.MaxSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d MB", config.MaxSize/(1<<20))
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !strings.Contains(config.AllowedTypes, ext) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", config.AllowedTypes)
	}

	if err := os.MkdirAll(config.UploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create uploads directory: %w", err)
	}

	filename := fmt.Sprintf("%s_%s%s", config.FileNamePrefix, uuid.New().String(), ext)
	filePath := filepath.Join(config.UploadDir, filename)

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	relativePath := strings.TrimPrefix(filePath, "./")
	return "/" + strings.ReplaceAll(relativePath, "\\", "/"), nil
}

func DeleteFile(filePath string) error {
	if filePath == "" || strings.HasPrefix(filePath, "http") {
		return nil
	}

	if !strings.HasPrefix(filePath, "./") && !strings.HasPrefix(filePath, "/") {
		filePath = "./" + filePath
	}
	if strings.HasPrefix(filePath, "/") {
		filePath = "." + filePath
	}

	return os.Remove(filePath)
}

func GetFormFile(r *http.Request, fieldName string) (multipart.File, *multipart.FileHeader, error) {
	file, header, err := r.FormFile(fieldName)
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, nil, nil
		}
		return nil, nil, fmt.Errorf("error reading file: %w", err)
	}
	return file, header, nil
}

func FileToBase64(file multipart.File, header *multipart.FileHeader) (string, error) {
	if header.Size > MaxFileSize {
		return "", fmt.Errorf("file size exceeds maximum allowed size of %d MB", MaxFileSize/(1<<20))
	}

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !strings.Contains(AllowedImageTypes, ext) {
		return "", fmt.Errorf("invalid file type. Allowed types: %s", AllowedImageTypes)
	}

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	mimeType := getMimeType(ext)

	base64String := base64.StdEncoding.EncodeToString(fileBytes)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64String), nil
}

func getMimeType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".webp":
		return "image/webp"
	default:
		return "image/jpeg"
	}
}
