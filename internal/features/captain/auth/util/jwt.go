package util

import (
	"time"

	"encore.app/internal/config"
	"encore.app/internal/pkg/errs"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	CaptainID int32  `json:"captain_id"`
	Phone     string `json:"phone"`
	TokenType string `json:"token_type"` // "access" or "refresh"
	jwt.RegisteredClaims
}

func GenerateAccessToken(captainID int32, phone string, cfg *config.Config) (string, error) {
	return generateToken(captainID, phone, "access", 15*time.Minute, cfg)
}

func GenerateRefreshToken(captainID int32, phone string, cfg *config.Config) (string, error) {
	return generateToken(captainID, phone, "refresh", 30*24*time.Hour, cfg)
}

func generateToken(captainID int32, phone, tokenType string, expiry time.Duration, cfg *config.Config) (string, error) {
	claims := Claims{
		CaptainID: captainID,
		Phone:     phone,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "encore.app-captain-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(cfg.JWTSecret))
	if err != nil {
		return "", errs.Internal(err, "failed to sign JWT token")
	}

	return tokenString, nil
}

func VerifyToken(tokenString string, cfg *config.Config) (int32, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.BadRequest("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", errs.Unauthorized("invalid or expired token")
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.CaptainID, claims.Phone, nil
	}

	return 0, "", errs.Unauthorized("invalid token")
}

func VerifyRefreshToken(tokenString string, cfg *config.Config) (int32, string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errs.BadRequest("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})

	if err != nil {
		return 0, "", errs.Unauthorized("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return 0, "", errs.Unauthorized("invalid refresh token")
	}

	if claims.TokenType != "refresh" {
		return 0, "", errs.Unauthorized("not a refresh token")
	}

	return claims.CaptainID, claims.Phone, nil
}
