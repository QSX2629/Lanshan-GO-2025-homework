package handler

import (
	"context"

	filemodel "github.com/aim/aim/internal/file/model"
	fileservice "github.com/aim/aim/internal/file/service"
	"github.com/aim/aim/internal/pkg/database"
)

// FileHandler provides handlers for file operations.
type FileHandler struct {
	svc *fileservice.FileService
}

// NewFileHandler creates a new FileHandler.
func NewFileHandler(db *database.DB, storageDir string) *FileHandler {
	return &FileHandler{svc: fileservice.NewFileService(db, storageDir)}
}

// SaveFileRequest is the input for saving a file record.
type SaveFileRequest struct {
	UserID   uint   `json:"user_id"`
	FileName string `json:"file_name"`
	FileURL  string `json:"file_url"`
	FileSize int64  `json:"file_size"`
	MimeType string `json:"mime_type"`
}

// SaveFile records uploaded file metadata.
func (h *FileHandler) SaveFile(_ context.Context, req *SaveFileRequest) (*filemodel.FileMeta, error) {
	return h.svc.SaveFile(req.UserID, req.FileName, req.FileURL, req.FileSize, req.MimeType)
}

// GetFile returns file metadata by ID.
func (h *FileHandler) GetFile(_ context.Context, fileID uint) (*filemodel.FileMeta, error) {
	return h.svc.GetFile(fileID)
}

// ListFiles returns files uploaded by a user.
func (h *FileHandler) ListFiles(_ context.Context, userID uint, offset, limit int) ([]filemodel.FileMeta, error) {
	return h.svc.ListFiles(userID, offset, limit)
}

// DeleteFile deletes a file record.
func (h *FileHandler) DeleteFile(_ context.Context, fileID, userID uint) error {
	return h.svc.DeleteFile(fileID, userID)
}
