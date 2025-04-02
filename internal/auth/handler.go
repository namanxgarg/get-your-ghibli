package auth

import (
    "get-your-ghibli/pkg/models"
    "get-your-ghibli/internal/db"
    "github.com/gin-gonic/gin"
    "net/http"
)

type EmailRequest struct {
    Email string `json:"email"`
}

type VerifyRequest struct {
    Email string `json:"email"`
    OTP   string `json:"otp"`
}

func RequestOTPHandler(c *gin.Context) {
    var req EmailRequest
    if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email"})
        return
    }

    otp, err := SaveOTP(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
        return
    }

    // For now, print OTP to console (simulate email sending)
    c.JSON(http.StatusOK, gin.H{"message": "OTP sent!", "debug_otp": otp})

    // Create user if not exists
    var user models.User
    result := db.DB.Where("email = ?", req.Email).First(&user)
    if result.RowsAffected == 0 {
        db.DB.Create(&models.User{Email: req.Email})
    }
}

func VerifyOTPHandler(c *gin.Context) {
    var req VerifyRequest
    if err := c.ShouldBindJSON(&req); err != nil || req.Email == "" || req.OTP == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
        return
    }

    if !VerifyOTP(req.Email, req.OTP) {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
        return
    }

    token, err := GenerateJWT(req.Email)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": token})
}
