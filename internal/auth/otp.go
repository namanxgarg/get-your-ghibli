package auth

import (
    "context"
    "crypto/rand"
    "fmt"
    "math/big"
    "time"

    "github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func InitRedis(redisURL string) {
    opt, err := redis.ParseURL(redisURL)
    if err != nil {
        panic(fmt.Sprintf("❌ Failed to parse Redis URL: %v", err))
    }

    rdb = redis.NewClient(opt)
    fmt.Println("✅ Connected to Redis")
}

func generateOTP() string {
    otp := ""
    for i := 0; i < 6; i++ {
        n, _ := rand.Int(rand.Reader, big.NewInt(10))
        otp += fmt.Sprintf("%d", n.Int64())
    }
    return otp
}

func SaveOTP(email string) (string, error) {
    otp := generateOTP()
    key := fmt.Sprintf("otp:%s", email)

    err := rdb.Set(ctx, key, otp, 5*time.Minute).Err()
    if err != nil {
        return "", err
    }

    return otp, nil
}

func VerifyOTP(email, input string) bool {
    key := fmt.Sprintf("otp:%s", email)
    storedOTP, err := rdb.Get(ctx, key).Result()
    if err != nil {
        return false
    }

    return storedOTP == input
}
