package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type CustomClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	UserId   int64  `json:"user_id"`
}
type JwtConfig struct {
	AccessSecret      string
	RefreshSecret     string
	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
}

func GenerateAccessToken(conf JwtConfig, userId int64, username string) (string, error) {
	claims := CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.AccessExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.AccessSecret))
}
func ParseAccessToken(conf JwtConfig, tokenString string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.AccessSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
func ParseRefreshToken(conf JwtConfig, tokenStr string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.RefreshSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid refresh token")
}
func GenerateRefreshToken(conf JwtConfig, userId int64, username string) (string, error) {
	claims := CustomClaims{
		UserId:   userId,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(conf.RefreshExpireTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(conf.RefreshSecret))
}
