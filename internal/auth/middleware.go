package auth

import (
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "net/http"
    "os"
    "strings"
)

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
            c.Abort()
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            // Ensure token is signed with correct method
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, jwt.ErrSignatureInvalid
            }
            return []byte(os.Getenv("JWT_SECRET")), nil
        })

        if err != nil || !token.Valid {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        claims, ok := token.Claims.(jwt.MapClaims)
        if !ok || claims["email"] == nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
            c.Abort()
            return
        }

        // Set user email in context for other routes to use
        c.Set("user_email", claims["email"])
        c.Next()
    }
}
