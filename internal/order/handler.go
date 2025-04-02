package order

import (
	"fmt"
	"get-your-ghibli/internal/db"
	"get-your-ghibli/pkg/models"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
)

type OrderRequest struct {
	GeneratedImageID uint   `json:"generated_image_id"`
	ShippingAddress  string `json:"shipping_address"`
	Note             string `json:"note"`
}

func PlaceOrderHandler(c *gin.Context) {
	email := c.GetString("user_email")

	var req OrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order data"})
		return
	}

	// Create DB record
	order := models.Order{
		UserEmail:        email,
		GeneratedImageID: req.GeneratedImageID,
		ShippingAddress:  req.ShippingAddress,
		Note:             req.Note,
		Status:           "pending",
	}
	db.DB.Create(&order)

	// Razorpay ₹2000 order
	razor := razorpay.NewClient(os.Getenv("RAZORPAY_KEY"), os.Getenv("RAZORPAY_SECRET"))
	amount := 2000 * 100 // ₹2000 in paise
	razorOrder, err := razor.Order.Create(map[string]interface{}{
		"amount":   amount,
		"currency": "INR",
		"receipt":  fmt.Sprintf("order_%d", order.ID),
	}, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create Razorpay order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order_id":      order.ID,
		"payment_order": razorOrder,
	})
}
