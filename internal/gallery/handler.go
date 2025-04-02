package gallery

import (
    "get-your-ghibli/internal/db"
    "get-your-ghibli/pkg/models"
    "github.com/gin-gonic/gin"
    "net/http"
)

func GalleryHandler(c *gin.Context) {
    email := c.GetString("user_email")

    var uploads []models.Upload
    db.DB.Where("user_email = ? AND status = ?", email, "ready").Find(&uploads)

    var result []models.UploadWithImages
    for _, upload := range uploads {
        var images []models.GeneratedImage
        db.DB.Where("upload_id = ?", upload.ID).Find(&images)

        result = append(result, models.UploadWithImages{
            Upload: upload,
            Images: images,
        })
    }

    c.JSON(http.StatusOK, gin.H{
        "gallery": result,
    })
}
