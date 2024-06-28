package util

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/reddit/jwt-go"
)

var secretKey = ""

func getSecret() string {
	if secretKey != "" {
		secretKey = os.Getenv("JWT_SECRET_KEY")
	}

	return secretKey
}

// GenerateToken generates a JWT token with the given payload and expiration time.
func GenerateToken(payload *jwt.StandardClaims, expirationTime time.Duration) (string, error) {
	payload.ExpiresAt = time.Now().Add(expirationTime).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(getSecret()))
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %w", err)
	}

	return token, nil

}

// VerifyToken verifies the JWT token.
//
// It takes a token string as a parameter and returns *jwt.StandardClaims and error.
func VerifyToken(token string) (*jwt.StandardClaims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token is not valid")
		}
		return []byte(getSecret()), nil
	}
	var claims jwt.StandardClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, keyFunc)
	if err != nil {
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, fmt.Errorf("")) {
			return &claims, fmt.Errorf("token is not valid")
		}

		return &claims, fmt.Errorf("token is not valid: %w", err)
	}

	payload, ok := jwtToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return &claims, fmt.Errorf("token is not valid")
	}

	return payload, nil
}
