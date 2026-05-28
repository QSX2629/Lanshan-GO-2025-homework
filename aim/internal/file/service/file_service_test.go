package service

import (
	"testing"

	filerepo "github.com/aim/aim/internal/file/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupFileService(t *testing.T) *FileService {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := filerepo.NewFileRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewFileService(db, "./upload")
}

func TestFileService_SaveAndGet(t *testing.T) {
	svc := setupFileService(t)

	f, err := svc.SaveFile(1, "photo.jpg", "/upload/photo.jpg", 1024, "image/jpeg")
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}
	if f.FileName != "photo.jpg" {
		t.Errorf("FileName = %q, want photo.jpg", f.FileName)
	}

	got, err := svc.GetFile(f.ID)
	if err != nil {
		t.Fatalf("GetFile() error = %v", err)
	}
	if got.FileURL != "/upload/photo.jpg" {
		t.Errorf("FileURL = %q", got.FileURL)
	}

	_, err = svc.GetFile(999)
	if err != ErrFileNotFound {
		t.Errorf("GetFile(999) error = %v, want ErrFileNotFound", err)
	}
}

func TestFileService_ListFiles(t *testing.T) {
	svc := setupFileService(t)
	svc.SaveFile(1, "a.jpg", "/a.jpg", 100, "image/jpeg")
	svc.SaveFile(1, "b.jpg", "/b.jpg", 200, "image/jpeg")

	files, err := svc.ListFiles(1, 0, 10)
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}
	if len(files) != 2 {
		t.Errorf("len(files) = %d, want 2", len(files))
	}

	// Test default limit.
	files2, err := svc.ListFiles(1, 0, 0)
	if err != nil {
		t.Fatalf("ListFiles(0 limit) error = %v", err)
	}
	if len(files2) != 2 {
		t.Errorf("len(files) with default limit = %d, want 2", len(files2))
	}
}

func TestFileService_Delete(t *testing.T) {
	svc := setupFileService(t)
	f, _ := svc.SaveFile(1, "tmp.txt", "/tmp.txt", 10, "text/plain")

	if err := svc.DeleteFile(f.ID, 1); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}

	if err := svc.DeleteFile(f.ID, 1); err != ErrFileNotFound {
		t.Errorf("delete again error = %v, want ErrFileNotFound", err)
	}
}

func TestFileService_DeleteNotOwner(t *testing.T) {
	svc := setupFileService(t)
	f, _ := svc.SaveFile(1, "tmp.txt", "/tmp.txt", 10, "text/plain")

	err := svc.DeleteFile(f.ID, 2)
	if err == nil {
		t.Error("expected permission denied error")
	}
}

func TestFileService_TooLarge(t *testing.T) {
	svc := setupFileService(t)

	_, err := svc.SaveFile(1, "big.bin", "/big.bin", maxFileSize+1, "application/octet-stream")
	if err != ErrFileTooLarge {
		t.Errorf("error = %v, want ErrFileTooLarge", err)
	}
}
