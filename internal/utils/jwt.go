package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

type Claims struct {
	UserID       string `json:"userId"`
	Role         string `json:"role"`
	IsApproved   bool   `json:"isApproved"`
	NeedsProfile bool   `json:"needsProfile"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID, role string, isApproved, needsProfile bool) string {
	claims := Claims{
		UserID:       userID,
		Role:         role,
		IsApproved:   isApproved,
		NeedsProfile: needsProfile,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		log.Println("Error signing JWT:", err)
		return ""
	}
	return tokenString
}

func ParseJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
