package models

import (
	"time"

	"gorm.io/gorm"
)

type Audio struct {
	ID          uint              `gorm:"primaryKey" json:"id"`
	UserID      uint              `gorm:"index;not null" json:"user_id"`
	User        User              `gorm:"foreignKey:UserID" json:"-"`
	Filename    string            `gorm:"size:255;not null" json:"filename"`
	FilePath    string            `gorm:"size:512;not null" json:"file_path"`      // Path in storage or local
	Duration    float64           `json:"duration"`                                // In seconds
	Status      string            `gorm:"size:50;default:'pending'" json:"status"` // pending, processing, completed, failed
	Language    string            `gorm:"size:10;default:'en'" json:"language"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	DeletedAt   gorm.DeletedAt    `gorm:"index" json:"-"`
	Transcripts []AudioTranscript `json:"transcripts,omitempty"`
}
