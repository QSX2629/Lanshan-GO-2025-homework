package model

import (
	"testing"
	"time"
)

func TestFileMeta_Image(t *testing.T) {
	f := FileMeta{
		UserID:   1,
		FileName: "photo.png",
		FileURL:  "/upload/photo.png",
		FileSize: 102400,
		MimeType: "image/png",
		Width:    1920,
		Height:   1080,
	}

	if f.MimeType != "image/png" {
		t.Errorf("MimeType = %q, want image/png", f.MimeType)
	}
	if f.Width != 1920 || f.Height != 1080 {
		t.Errorf("dimensions = %dx%d, want 1920x1080", f.Width, f.Height)
	}
}

func TestFileMeta_Voice(t *testing.T) {
	f := FileMeta{
		UserID:   1,
		FileName: "voice.mp3",
		FileURL:  "/upload/voice.mp3",
		FileSize: 51200,
		MimeType: "audio/mpeg",
		Duration: 60,
	}

	if f.Duration != 60 {
		t.Errorf("Duration = %d, want 60", f.Duration)
	}
}

func TestFileMeta_General(t *testing.T) {
	f := FileMeta{
		ID:        1,
		UserID:    1,
		FileName:  "doc.pdf",
		FileURL:   "/upload/doc.pdf",
		FileSize:  204800,
		MimeType:  "application/pdf",
		CreatedAt: time.Now(),
	}

	if f.FileName != "doc.pdf" {
		t.Errorf("FileName = %q, want doc.pdf", f.FileName)
	}
	if f.FileSize != 204800 {
		t.Errorf("FileSize = %d, want 204800", f.FileSize)
	}
}
