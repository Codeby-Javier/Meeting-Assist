package models

import (
	"time"

	"gorm.io/gorm"
)

type Document struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"index;not null" json:"user_id"`
	User      User           `gorm:"foreignKey:UserID" json:"-"`
	Filename  string         `gorm:"size:255;not null" json:"filename"`
	FilePath  string         `gorm:"size:512;not null" json:"file_path"`
	Type      string         `gorm:"size:50" json:"type"` // pdf, image
	Status    string         `gorm:"size:50;default:'pending'" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Texts     []DocumentText `json:"texts,omitempty"`
}
