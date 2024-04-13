package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"time"

	"github.com/Alladan04/avito_test/internal/models"
	"github.com/golang-jwt/jwt/v5"
)

func GetHash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashInBytes := hasher.Sum(nil)
	return hex.EncodeToString(hashInBytes)
}

func GenToken(user models.User, lifeTime time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"usr": user.Username,
		"exp": time.Now().Add(lifeTime).Unix(),
	})
	if user.IsAdmin {
		return token.SignedString([]byte(os.Getenv("JWT_ADMIN_SECRET")))
	}
	return token.SignedString([]byte(os.Getenv("JWT_USER_SECRET")))
}
