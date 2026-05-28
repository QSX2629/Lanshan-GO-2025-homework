package handler

import (
	"testing"

	filerepo "github.com/aim/aim/internal/file/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupFileHandler(t *testing.T) *FileHandler {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := filerepo.NewFileRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewFileHandler(db, "./upload")
}

func TestFileHandler_SaveAndGet(t *testing.T) {
	h := setupFileHandler(t)

	f, err := h.SaveFile(nil, &SaveFileRequest{
		UserID: 1, FileName: "photo.jpg", FileURL: "/upload/photo.jpg",
		FileSize: 1024, MimeType: "image/jpeg",
	})
	if err != nil {
		t.Fatalf("SaveFile() error = %v", err)
	}
	if f.FileName != "photo.jpg" {
		t.Errorf("FileName = %q, want photo.jpg", f.FileName)
	}

	got, err := h.GetFile(nil, f.ID)
	if err != nil {
		t.Fatalf("GetFile() error = %v", err)
	}
	if got.FileSize != 1024 {
		t.Errorf("FileSize = %d, want 1024", got.FileSize)
	}
}

func TestFileHandler_ListFiles(t *testing.T) {
	h := setupFileHandler(t)
	h.SaveFile(nil, &SaveFileRequest{UserID: 1, FileName: "a.jpg", FileURL: "/a.jpg", FileSize: 100, MimeType: "image/jpeg"})
	h.SaveFile(nil, &SaveFileRequest{UserID: 1, FileName: "b.jpg", FileURL: "/b.jpg", FileSize: 200, MimeType: "image/jpeg"})

	files, err := h.ListFiles(nil, 1, 0, 10)
	if err != nil {
		t.Fatalf("ListFiles() error = %v", err)
	}
	if len(files) != 2 {
		t.Errorf("len(files) = %d, want 2", len(files))
	}
}

func TestFileHandler_Delete(t *testing.T) {
	h := setupFileHandler(t)
	f, _ := h.SaveFile(nil, &SaveFileRequest{UserID: 1, FileName: "tmp.txt", FileURL: "/tmp.txt", FileSize: 10, MimeType: "text/plain"})

	if err := h.DeleteFile(nil, f.ID, 1); err != nil {
		t.Fatalf("DeleteFile() error = %v", err)
	}
}
