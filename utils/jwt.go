package utils

import (
	"volunteer-system-backend/config"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"sync"
	"time"
)

var (
	jwtKey []byte
	once   sync.Once
)

// 初始化 JWT 密钥
func initJwtKey() {
	once.Do(func() {
		jwtKey = []byte(GenerateMD5(config.ProjectConfig.Volunteer.TwtKey))
	})
}

// Claims 结构体
type Claims struct {
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	LoginTime time.Time `json:"loginTime"`
	jwt.RegisteredClaims
}

// GenerateJWT 生成 JWT
func GenerateJWT(email, nickname string, loginTime time.Time) (string, error) {
	initJwtKey() // 确保 jwtKey 已初始化

	// 设置过期时间
	var expMinutes int64 = 1 // 默认1小时
	if config.ProjectConfig.Volunteer.JwtExpiry > 0 {
		expMinutes = int64(config.ProjectConfig.Volunteer.JwtExpiry)
	}

	// 计算过期时间
	expirationTime := time.Now().Add(time.Duration(expMinutes) * time.Hour)

	claims := &Claims{
		Email:     email,
		Nickname:  nickname,
		LoginTime: loginTime,
		RegisteredClaims: jwt.RegisteredClaims{
			// 添加更多标准字段来确保token的唯一性和时效性
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),           // 添加签发时间
			NotBefore: jwt.NewNumericDate(time.Now()),           // 添加生效时间
			ID:        GenerateMD5(email + time.Now().String()), // 添加唯一标识符
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT 验证 JWT
func ParseJWT(tokenStr string) (*Claims, error) {
	initJwtKey() // 确保 jwtKey 已初始化

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("意外的签名方法")
		}
		return jwtKey, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, errors.New("token 已过期")
		}
		return nil, errors.New("token 解析失败: " + err.Error())
	}

	// 检查token是否有效
	if !token.Valid {
		return nil, errors.New("token 无效")
	}

	// 严格检查过期时间
	if claims.ExpiresAt != nil {
		if time.Until(claims.ExpiresAt.Time) <= 0 {
			return nil, errors.New("token 已过期")
		}
	}

	return claims, nil
}

// GenerateMD5 生成字符串的 MD5 哈希
func GenerateMD5(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
