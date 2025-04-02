package models

type UploadWithImages struct {
    Upload
    Images []GeneratedImage `json:"images"`
}
