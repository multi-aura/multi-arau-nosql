package jwt

import (
	config "multiaura/internal/configs/dev"
	"multiaura/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["userID"] = user.ID
	claims["fullname"] = user.FullName
	claims["email"] = user.Email
	claims["phone"] = user.PhoneNumber
	claims["isAdmin"] = user.IsAdmin
	claims["isActive"] = user.IsActive
	claims["exp"] = time.Now().Add(time.Hour * 168).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	cfg, err := config.Instance()
	if err != nil {
		return "", err
	}

	return token.SignedString([]byte(cfg.GetSecretKey()))
}