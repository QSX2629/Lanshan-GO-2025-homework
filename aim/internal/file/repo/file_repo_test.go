package repo

import (
	"testing"

	filemodel "github.com/aim/aim/internal/file/model"
	"github.com/aim/aim/internal/pkg/database"
)

func setupFileRepo(t *testing.T) *FileRepo {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := NewFileRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return repo
}

func TestFileRepo_CreateAndFind(t *testing.T) {
	repo := setupFileRepo(t)

	f := &filemodel.FileMeta{
		UserID:   1,
		FileName: "photo.png",
		FileURL:  "/upload/photo.png",
		FileSize: 1024,
		MimeType: "image/png",
		Width:    800,
		Height:   600,
	}
	if err := repo.Create(f); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	found, err := repo.FindByID(f.ID)
	if err != nil {
		t.Fatalf("FindByID() error = %v", err)
	}
	if found.FileName != "photo.png" {
		t.Errorf("FileName = %q, want photo.png", found.FileName)
	}

	files, err := repo.FindByUserID(1, 0, 10)
	if err != nil {
		t.Fatalf("FindByUserID() error = %v", err)
	}
	if len(files) != 1 {
		t.Errorf("len(files) = %d, want 1", len(files))
	}
}

func TestFileRepo_Delete(t *testing.T) {
	repo := setupFileRepo(t)

	f := &filemodel.FileMeta{UserID: 1, FileName: "tmp.txt", FileURL: "/tmp.txt", FileSize: 10}
	repo.Create(f)

	if err := repo.Delete(f.ID); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	_, err := repo.FindByID(f.ID)
	if err == nil {
		t.Error("expected error for deleted file")
	}
}
