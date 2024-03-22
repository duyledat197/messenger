package util

import (
	"errors"
	"fmt"
	"time"

	"github.com/reddit/jwt-go"
)

// JWTAuthenticator is representation of [Authenticator] engine that implement using JWT.
type JWTAuthenticator struct {
	secretKey string
}

// NewJWTAuthenticator creates a new JWTAuthenticator.
//
// It takes a secret key as a string parameter and returns a JWTAuthenticator and an error.
func NewJWTAuthenticator(secretKey string) (JWTAuthenticator, error) {
	return JWTAuthenticator{
		secretKey,
	}, nil
}

// Generate generates a JWT token with the given payload and expiration time.
//
// Parameters:
// - payload: The payload containing the claims for the JWT token.
// - expirationTime: The duration for which the token should be valid.
//
// Returns:
// - string: The generated JWT token.
// - error: An error if the token generation fails.
func (a *JWTAuthenticator) Generate(payload *jwt.StandardClaims, expirationTime time.Duration) (string, error) {
	payload.ExpiresAt = time.Now().Add(expirationTime).Unix()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	token, err := jwtToken.SignedString([]byte(a.secretKey))
	if err != nil {
		return "", fmt.Errorf("unable to generate token: %w", err)
	}

	return token, nil

}

// Verify verifies the JWT token.
//
// It takes a token string as a parameter and returns *jwt.StandardClaims and error.
func (a *JWTAuthenticator) Verify(token string) (*jwt.StandardClaims, error) {
	keyFunc := func(token *jwt.Token) (any, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("token is not valid")
		}
		return []byte(a.secretKey), nil
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
