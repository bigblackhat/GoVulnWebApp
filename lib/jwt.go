package lib

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

// 定义JWT声明结构体，需要与生成JWT时的结构匹配
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成包含username和过期时间的JWT令牌
// username: 要包含在JWT中的用户名
// secretKey: 用于签名的密钥
// expiration: JWT的有效时长（如24*time.Hour）
func GenerateJWT(username, secretKey string, expiration time.Duration) (string, error) {
	// 确保密钥长度安全（至少32字节）
	if len(secretKey) < 32 {
		return "", fmt.Errorf("secret key must be at least 32 bytes long for security")
	}

	// 创建声明
	claims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)), // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                 // 签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),                 // 生效时间
			Issuer:    "your-app-name",                                // 签发者
			Subject:   "user-access",                                  // 主题
			// ID: 已移除硬编码的ID，应使用随机值或基于用户生成的ID
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名令牌
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("生成JWT失败: %w", err)
	}

	return tokenString, nil
}

// ParseJWT 解析JWT令牌并返回username和过期时间
func ParseJWT(tokenString, secretKey string) (username string, expiration time.Time, err error) {
	// 确保密钥长度安全
	if len(secretKey) < 32 {
		return "", time.Time{}, fmt.Errorf("secret key must be at least 32 bytes long for security")
	}

	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to parse JWT: %w", err)
	}

	// 验证令牌是否有效
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return "", time.Time{}, fmt.Errorf("invalid JWT token")
	}

	// 检查过期时间是否存在
	if claims.ExpiresAt == nil {
		return "", time.Time{}, fmt.Errorf("token is missing expiration claim")
	}

	// 返回username和过期时间
	return claims.Username, claims.ExpiresAt.Time, nil
}
