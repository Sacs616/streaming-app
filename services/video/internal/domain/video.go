package domain

import "time"

type VideoStatus string

const (
	StatusPending     VideoStatus = "pending"
	StatusTranscoding VideoStatus = "transcoding"
	StatusReady       VideoStatus = "ready"
	StatusFailed      VideoStatus = "failed"
)

type Video struct {
	ID           uint64      `gorm:"primaryKey"`
	Title        string      `gorm:"not null"`
	Description  string      `gorm:"type:text"`
	UploadedBy   uint64      `gorm:"not null"`
	Duration     int         `gorm:"default:0"` // in seconds
	ThumbnailURL string      `gorm:"column:thumbnail_url"`
	HLSManifest  string      `gorm:"column:hls_manifest"`
	Status       VideoStatus `gorm:"default:'pending'"`
	CategoryID   uint64      `gorm:"not null"`
	ViewCount    int         `gorm:"default:0"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func (Video) TableName() string {
	return "videos"
}
