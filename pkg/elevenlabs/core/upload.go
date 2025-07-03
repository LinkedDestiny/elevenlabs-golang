package core

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// FileUpload represents a file to be uploaded
type FileUpload struct {
	FieldName string
	FileName  string
	Content   io.Reader
	MimeType  string
}

// CreateMultipartRequest creates an HTTP request with multipart form data
func CreateMultipartRequest(files []FileUpload, fields map[string]string) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add form fields
	for key, value := range fields {
		if err := writer.WriteField(key, value); err != nil {
			writer.Close()
			return nil, fmt.Errorf("failed to write field %s: %w", key, err)
		}
	}

	// Add files
	for _, file := range files {
		var part io.Writer
		var err error

		if file.FileName != "" {
			part, err = writer.CreateFormFile(file.FieldName, file.FileName)
		} else {
			part, err = writer.CreateFormField(file.FieldName)
		}

		if err != nil {
			writer.Close()
			return nil, fmt.Errorf("failed to create form file for field %s: %w", file.FieldName, err)
		}

		if _, err := io.Copy(part, file.Content); err != nil {
			writer.Close()
			return nil, fmt.Errorf("failed to copy file content for field %s: %w", file.FieldName, err)
		}
	}

	// Close the writer to finalize the form
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("failed to close multipart writer: %w", err)
	}

	// Create the request
	req, err := http.NewRequest("POST", "", body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, nil
}

// FileUploadFromPath creates a FileUpload from a file path
func FileUploadFromPath(fieldName, filePath string) (*FileUpload, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	fileName := filepath.Base(filePath)
	mimeType := getMimeTypeFromExtension(filepath.Ext(filePath))

	return &FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Content:   file,
		MimeType:  mimeType,
	}, nil
}

// FileUploadFromBytes creates a FileUpload from byte data
func FileUploadFromBytes(fieldName, fileName string, data []byte) *FileUpload {
	mimeType := getMimeTypeFromExtension(filepath.Ext(fileName))

	return &FileUpload{
		FieldName: fieldName,
		FileName:  fileName,
		Content:   bytes.NewReader(data),
		MimeType:  mimeType,
	}
}

// getMimeTypeFromExtension returns the MIME type for a file extension
func getMimeTypeFromExtension(ext string) string {
	switch ext {
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".flac":
		return "audio/flac"
	case ".m4a":
		return "audio/mp4"
	case ".ogg":
		return "audio/ogg"
	case ".aac":
		return "audio/aac"
	case ".txt":
		return "text/plain"
	case ".pdf":
		return "application/pdf"
	case ".json":
		return "application/json"
	default:
		return "application/octet-stream"
	}
}

// ValidateAudioFile checks if a file is a valid audio file
func ValidateAudioFile(filePath string) error {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Check file extension
	ext := filepath.Ext(filePath)
	validExtensions := []string{".mp3", ".wav", ".flac", ".m4a", ".ogg", ".aac"}

	for _, validExt := range validExtensions {
		if ext == validExt {
			return nil
		}
	}

	return fmt.Errorf("unsupported audio file format: %s", ext)
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filePath string) (int64, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to get file info: %w", err)
	}

	return info.Size(), nil
}
