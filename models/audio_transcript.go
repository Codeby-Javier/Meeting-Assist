package models

import (
	"time"
)

type AudioTranscript struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AudioID   uint      `gorm:"index;not null" json:"audio_id"`
	Text      string    `gorm:"type:text" json:"text"`
	StartTime float64   `json:"start_time"`
	EndTime   float64   `json:"end_time"`
	Speaker   string    `gorm:"size:100" json:"speaker"`
	CreatedAt time.Time `json:"created_at"`
}
