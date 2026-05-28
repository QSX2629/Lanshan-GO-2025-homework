package repo

import (
	"github.com/aim/aim/internal/file/model"
	"github.com/aim/aim/internal/pkg/database"
)

// FileRepo handles file metadata data access.
type FileRepo struct {
	db *database.DB
}

// NewFileRepo creates a new FileRepo.
func NewFileRepo(db *database.DB) *FileRepo {
	return &FileRepo{db: db}
}

// AutoMigrate creates the file metadata table.
func (r *FileRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&model.FileMeta{})
}

// Create inserts a file metadata record.
func (r *FileRepo) Create(f *model.FileMeta) error {
	return r.db.Create(f).Error
}

// FindByID looks up a file record by ID.
func (r *FileRepo) FindByID(id uint) (*model.FileMeta, error) {
	var f model.FileMeta
	err := r.db.First(&f, id).Error
	if err != nil {
		return nil, err
	}
	return &f, nil
}

// FindByUserID returns all files uploaded by a user.
func (r *FileRepo) FindByUserID(userID uint, offset, limit int) ([]model.FileMeta, error) {
	var files []model.FileMeta
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&files).Error
	return files, err
}

// Delete removes a file record.
func (r *FileRepo) Delete(id uint) error {
	return r.db.Delete(&model.FileMeta{}, id).Error
}
