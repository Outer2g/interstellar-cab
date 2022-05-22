package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type claims struct {
	Email string `json:"email"`
	Vip   bool   `json:"vip"`
	jwt.StandardClaims
}

var JWT_SECRET string

func GenerateJwtToken(email string, vip bool) (string, error) {
	key := []byte(JWT_SECRET)

	expirationTime := time.Now().Add(7 * 24 * 60 * time.Minute)
	claims := &claims{
		Email: email,
		Vip:   vip,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	UnsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := UnsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}

func VerifyJwtToken(strToken string) (*claims, error) {
	key := []byte(JWT_SECRET)

	claims := &claims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, fmt.Errorf("ERROR invalid token signature")
		}
	}

	if !token.Valid {
		return claims, fmt.Errorf("ERROR invalid token")
	}

	return claims, nil
}
