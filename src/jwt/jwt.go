package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var jwtSecret = []byte("luca_go_hello_go_secret_key")

const (
	Issuer          = "hello-go"
	TokenExpiration = 24 * 7 * time.Hour // 7 days
)

// TokenClaims 自定义 claims 结构
type TokenClaims struct {
	UserId   int64    `json:"userId"`
	Username string   `json:"username"`
	Roles    []string `json:"roles,omitempty"`
	jwt.RegisteredClaims
}

// GenerateToken 生成 JWT token
func GenerateToken(userId int64, username string, roles []string) (string, error) {
	now := time.Now()
	claims := TokenClaims{
		UserId:   userId,
		Username: username,
		Roles:    roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    Issuer,                                       // 签发者
			Subject:   username,                                     // 主题
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenExpiration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(now),                      // 签发时间
			NotBefore: jwt.NewNumericDate(now),                      // 生效时间
			ID:        uuid.New().String(),                          // 唯一ID，防重放
		},
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

// GetClaimsFromToken 从 token 提取完整 claims
func GetClaimsFromToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, jwt.ErrTokenInvalidClaims
}

// GetUserIdFromToken 提取 userId 从 token
func GetUserIdFromToken(tokenString string) (int64, error) {
	claims, err := GetClaimsFromToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserId, nil
}

// IsTokenExpired 判断 token 是否过期
func IsTokenExpired(tokenString string) (bool, error) {
	claims, err := GetClaimsFromToken(tokenString)
	if err != nil {
		return true, err
	}
	if claims.ExpiresAt != nil {
		return time.Now().After(claims.ExpiresAt.Time), nil
	}
	return true, nil
}

// HasRole 检查用户是否拥有指定角色
func HasRole(tokenString string, role string) (bool, error) {
	claims, err := GetClaimsFromToken(tokenString)
	if err != nil {
		return false, err
	}
	for _, r := range claims.Roles {
		if r == role {
			return true, nil
		}
	}
	return false, nil
}
