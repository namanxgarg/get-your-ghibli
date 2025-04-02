package upload

import (
	"fmt"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/pkg/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

func UploadHandler(c *gin.Context) {
	email := c.GetString("user_email")

	// Parse uploaded file
	file, err := c.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file"})
		return
	}

	// For now: save file locally (we’ll do S3/Cloudinary next)
	filePath := fmt.Sprintf("uploads/%s_%s", email, file.Filename)
	err = c.SaveUploadedFile(file, filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Create DB record
	upload := models.Upload{
		UserEmail: email,
		FileName:  file.Filename,
		UploadURL: filePath, // this will be replaced with public URL later
		Status:    "pending",
	}
	db.DB.Create(&upload)

	// Create Razorpay ₹100 order
	razor := razorpay.NewClient(os.Getenv("RAZORPAY_KEY"), os.Getenv("RAZORPAY_SECRET"))
	amount := 100 * 100 // ₹100 in paise
	orderData := map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  fmt.Sprintf("upload_%d", upload.ID),
	}
	order, err := razor.Order.Create(orderData, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"upload_id":     upload.ID,
		"photo_name":    upload.FileName,
		"payment_order": order,
	})
}
