package middleware

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID         string `json:"user_id"`
	Role           string `json:"role"`
	OrganizationID string `json:"organization_id,omitempty"`
	jwt.RegisteredClaims
}

type JWTValidator struct {
	secretKey []byte
}

func NewJWTValidator(secretKey string) *JWTValidator {
	return &JWTValidator{
		secretKey: []byte(secretKey),
	}
}

func (jv *JWTValidator) ValidateToken(tokenString string) (*JWTClaims, error) {

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Printf("❌ Unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jv.secretKey, nil
	})

	if err != nil {
		log.Printf("❌ Token parsing error: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		// Debug current time vs expiration
		now := time.Now()

		// Check if token is expired
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(now) {
			log.Printf("❌ Token is expired")
			return nil, errors.New("token expired")
		}

		log.Printf("✅ Token is valid for user: %s, role: %s", claims.UserID, claims.Role)
		return claims, nil
	}

	log.Printf("❌ Invalid token or claims")
	return nil, errors.New("invalid token")
}

func (jv *JWTValidator) GenerateToken(userID, role, organizationID string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // 24 hours
	issuedAt := time.Now()

	claims := &JWTClaims{
		UserID:         userID,
		Role:           role,
		OrganizationID: organizationID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(issuedAt),
		},
	}

	// Debug logging
	log.Printf("🔐 Generating token for user: %s", userID)
	log.Printf("🔐 Role: %s", role)
	log.Printf("🔐 Issued at: %v", issuedAt)
	log.Printf("🔐 Expires at: %v", expirationTime)
	log.Printf("🔑 Secret key length: %d", len(jv.secretKey))
	log.Printf("🔑 Secret key: %s", jv.secretKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jv.secretKey)

	if err != nil {
		log.Printf("❌ Token generation error: %v", err)
		return "", err
	}

	log.Printf("✅ Token generated successfully")
	return tokenString, nil
}
