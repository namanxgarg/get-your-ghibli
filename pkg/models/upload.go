package models

import "time"

type Upload struct {
	ID         uint      `gorm:"primaryKey"`
	UserEmail  string    `gorm:"index"`
	FileName   string
	UploadURL  string
	Status     string    // pending, paid, generating, ready
	CreatedAt  time.Time
}
