package model

import "time"

// FileMeta stores metadata for an uploaded file.
type FileMeta struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	FileName  string    `gorm:"size:256;not null" json:"file_name"`
	FileURL   string    `gorm:"size:512;not null" json:"file_url"`
	FileSize  int64     `gorm:"not null" json:"file_size"`
	MimeType  string    `gorm:"size:128" json:"mime_type"`
	Width     int       `json:"width,omitempty"`    // image/video
	Height    int       `json:"height,omitempty"`   // image/video
	Duration  int       `json:"duration,omitempty"` // voice/video seconds
	CreatedAt time.Time `json:"created_at"`
}
