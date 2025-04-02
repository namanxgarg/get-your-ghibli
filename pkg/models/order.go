package models

import "time"

type Order struct {
	ID                uint      `gorm:"primaryKey"`
	UserEmail         string    `gorm:"index"`
	GeneratedImageID  uint
	ShippingAddress   string
	Note              string
	Status            string    // pending, paid, shipped
	CreatedAt         time.Time
}
