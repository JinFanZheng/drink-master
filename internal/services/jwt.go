package services

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ddteam/drink-master/internal/middleware"
	"github.com/ddteam/drink-master/internal/models"
)

// JWTService handles JWT token operations
type JWTService struct {
	secret []byte
	expiry time.Duration
}

// NewJWTService creates a new JWT service
func NewJWTService() *JWTService {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// #nosec G101 - This is a fallback for development only
		secret = "default_jwt_secret_change_this_in_production"
	}

	// Get expiry hours from environment, default to 24 hours
	expiryHours := 24
	if expiryStr := os.Getenv("JWT_EXPIRES_HOURS"); expiryStr != "" {
		if hours, err := strconv.Atoi(expiryStr); err == nil {
			expiryHours = hours
		}
	}

	return &JWTService{
		secret: []byte(secret),
		expiry: time.Duration(expiryHours) * time.Hour,
	}
}

// GenerateToken generates a JWT token for a member
func (j *JWTService) GenerateToken(member *models.Member) (string, error) {
	now := time.Now()
	claims := &middleware.JWTClaims{
		MemberID: member.ID,
		Role:     fmt.Sprintf("%d", member.Role), // Convert int to string
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(j.expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	// Set machine owner ID if the member is an owner
	if member.MachineOwnerId != nil {
		claims.MachineOwnerID = *member.MachineOwnerId
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// ValidateToken validates a JWT token and returns claims
func (j *JWTService) ValidateToken(tokenString string) (*middleware.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &middleware.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrTokenNotValidYet
	}

	claims, ok := token.Claims.(*middleware.JWTClaims)
	if !ok {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}
