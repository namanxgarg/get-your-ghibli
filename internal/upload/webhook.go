package upload

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/pkg/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
	"get-your-ghibli/internal/queue"
)

type RazorpayWebhook struct {
	Event   string `json:"event"`
	Payload struct {
		Payment struct {
			Entity struct {
				OrderID string `json:"order_id"`
				Status  string `json:"status"`
			} `json:"entity"`
		} `json:"payment"`
	} `json:"payload"`
}

// verifySignature validates Razorpay's webhook signature
func verifySignature(body []byte, headerSig string, secret string) bool {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write(body)
	computed := hex.EncodeToString(h.Sum(nil))
	return computed == headerSig
}

func RazorpayWebhookHandler(c *gin.Context) {
	secret := os.Getenv("RAZORPAY_WEBHOOK_SECRET")
	sig := c.GetHeader("X-Razorpay-Signature")

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unable to read body"})
		return
	}

	if !verifySignature(body, sig, secret) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid signature"})
		return
	}

	var event RazorpayWebhook
	if err := json.Unmarshal(body, &event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payload"})
		return
	}

	if event.Event == "payment.captured" {
		orderID := event.Payload.Payment.Entity.OrderID
		// orderID is like "order_upload_3"
		var uploadID int
		n, err := fmt.Sscanf(orderID, "upload_%d", &uploadID)
		if err != nil || n != 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
			return
		}

		// update upload status
		db.DB.Model(&models.Upload{}).
			Where("id = ?", uploadID).
			Update("status", "paid")

		fmt.Printf("✅ Upload #%d marked as PAID\n", uploadID)

		upload := models.Upload{}
		db.DB.First(&upload, uploadID)

		err = queue.SendTask(upload.ID, upload.UploadURL)
		if err != nil {
			fmt.Printf("❌ Failed to queue task: %v\n", err)
		} else {
			fmt.Printf("📤 Queued upload #%d for generation\n", upload.ID)
		}
	}

	

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
