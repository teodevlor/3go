package jwt

import (
	"errors"
	"fmt"
	"go-structure/global"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

type AdminClaims struct {
	AdminID string `json:"admin_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(accountID uuid.UUID) (string, error) {
	secret := global.Config.JwtConfig.Secret

	ttlStr := global.Config.JwtConfig.AccessTokenTtl
	if ttlStr == "" {
		ttlStr = "15m"
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 15 * time.Minute
	}

	claims := Claims{
		AccountID: accountID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateAdminAccessToken(adminID uuid.UUID) (string, error) {
	secret := global.Config.JwtConfig.Secret

	ttlStr := global.Config.JwtConfig.AccessTokenTtl
	if ttlStr == "" {
		ttlStr = "15m"
	}

	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		ttl = 15 * time.Minute
	}

	claims := AdminClaims{
		AdminID: adminID.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseAccessToken(tokenStr string) (uuid.UUID, error) {
	secret := global.Config.JwtConfig.Secret
	if secret == "" {
		secret = "dev-secret-key"
	}

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	return uuid.Parse(claims.AccountID)
}

func ParseAdminAccessToken(tokenStr string) (uuid.UUID, error) {
	secret := global.Config.JwtConfig.Secret
	if secret == "" {
		secret = "dev-secret-key"
	}

	token, err := jwt.ParseWithClaims(tokenStr, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	claims, ok := token.Claims.(*AdminClaims)
	if !ok || !token.Valid {
		return uuid.Nil, errors.New("invalid token")
	}

	return uuid.Parse(claims.AdminID)
}
