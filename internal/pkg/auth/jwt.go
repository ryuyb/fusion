package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ryuyb/fusion/internal/infrastructure/config"
	errors2 "github.com/ryuyb/fusion/internal/pkg/errors"
)

type UserClaims struct {
	UserID   int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secret     []byte
	expiration time.Duration
}

func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret:     []byte(cfg.Jwt.Secret),
		expiration: cfg.Jwt.Expiration,
	}
}

func (j *JWTManager) GenerateToken(userID int64, username, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.expiration)
	claims := UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "fusion",
			Subject:   fmt.Sprintf("user_%d", userID),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	signedString, err := token.SignedString(j.secret)
	return signedString, expiresAt, err
}

func (j *JWTManager) ValidateToken(tokenStr string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(token *jwt.Token) (any, error) {
		return j.secret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors2.Unauthorized("Token is expired")
		}
		return nil, errors2.Unauthorized("Token is invalid").Wrap(err)
	}
	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors2.Unauthorized("Token is invalid")
}
