package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("luca_go_hello_go_secret_key")

// 生成 JWT token
func GenerateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"exp":    time.Now().Add(24 * 7 * time.Hour).Unix(), // 7 days
		"iat":    time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// 解析 JWT token
func ParseToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
}

// 提取 userId 从 token
func GetUserIdFromToken(tokenString string) (uint, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userIdFloat, ok := claims["userId"].(float64); ok {
			return uint(userIdFloat), nil
		}
	}
	return 0, nil
}

// 判断 token 是否过期
func IsTokenExpired(tokenString string) (bool, error) {
	token, err := ParseToken(tokenString)
	if err != nil {
		return false, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if expFloat, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(expFloat), 0)
			return time.Now().After(expirationTime), nil
		}
	}
	return false, nil
}
