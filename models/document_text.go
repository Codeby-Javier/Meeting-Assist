package models

import (
	"time"
)

type DocumentText struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	DocumentID uint      `gorm:"index;not null" json:"document_id"`
	Text       string    `gorm:"type:text" json:"text"`
	PageNumber int       `json:"page_number"`
	CreatedAt  time.Time `json:"created_at"`
}
