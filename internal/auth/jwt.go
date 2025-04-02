package auth

import (
    "github.com/golang-jwt/jwt/v5"
    "os"
    "time"
)

func GenerateJWT(email string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "exp":   time.Now().Add(24 * time.Hour).Unix(), // 1-day expiry
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
