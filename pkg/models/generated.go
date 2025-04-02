package models

type GeneratedImage struct {
    ID       uint   `gorm:"primaryKey"`
    UploadID uint
    ImageURL string
}
