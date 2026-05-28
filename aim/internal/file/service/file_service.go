package service

import (
	"errors"
	"time"

	filemodel "github.com/aim/aim/internal/file/model"
	filerepo "github.com/aim/aim/internal/file/repo"
	"github.com/aim/aim/internal/pkg/database"
)

var (
	ErrFileNotFound = errors.New("file not found")
	ErrFileTooLarge = errors.New("file too large")
)

const maxFileSize = 100 * 1024 * 1024 // 100 MB

// FileService handles file upload and management business logic.
type FileService struct {
	repo       *filerepo.FileRepo
	storageDir string // local storage path
}

// NewFileService creates a new FileService.
func NewFileService(db *database.DB, storageDir string) *FileService {
	return &FileService{repo: filerepo.NewFileRepo(db), storageDir: storageDir}
}

// SaveFile records file metadata after an upload.
func (s *FileService) SaveFile(userID uint, fileName, fileURL string, fileSize int64, mimeType string) (*filemodel.FileMeta, error) {
	if fileSize > maxFileSize {
		return nil, ErrFileTooLarge
	}

	f := &filemodel.FileMeta{
		UserID:    userID,
		FileName:  fileName,
		FileURL:   fileURL,
		FileSize:  fileSize,
		MimeType:  mimeType,
		CreatedAt: time.Now(),
	}
	if err := s.repo.Create(f); err != nil {
		return nil, err
	}
	return f, nil
}

// GetFile returns file metadata by ID.
func (s *FileService) GetFile(fileID uint) (*filemodel.FileMeta, error) {
	f, err := s.repo.FindByID(fileID)
	if err != nil {
		return nil, ErrFileNotFound
	}
	return f, nil
}

// ListFiles returns files uploaded by a user.
func (s *FileService) ListFiles(userID uint, offset, limit int) ([]filemodel.FileMeta, error) {
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	return s.repo.FindByUserID(userID, offset, limit)
}

// DeleteFile removes a file record.
func (s *FileService) DeleteFile(fileID, userID uint) error {
	f, err := s.repo.FindByID(fileID)
	if err != nil {
		return ErrFileNotFound
	}
	if f.UserID != userID {
		return errors.New("permission denied: not the file owner")
	}
	return s.repo.Delete(fileID)
}
