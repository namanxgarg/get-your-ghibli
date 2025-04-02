package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Email     string    `gorm:"uniqueIndex" json:"email"`
    CreatedAt time.Time `json:"created_at"`
}
