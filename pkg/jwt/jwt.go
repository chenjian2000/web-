package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const TokenExpireDuration = time.Hour * 2

var mySecret = []byte("niko")

// 定义 payload （claims）
type MyClaims struct {
	UserID               uint64 `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // jwt.StandardClaims 已经被弃用
}

// GenToken 生成 JWT token
func GenToken(user_id uint64, username string) (string, error) {
	// 生成 payload 部分
	claims := &MyClaims{
		UserID:   user_id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExpireDuration)), // 过期时间
			Issuer:    "niko-goweb",                                            // 签发人
		},
	}
	// 使用 HS256（HMAC SHA-256）签名算法创建 JWT 对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用密钥签名并生成最终的 JWT 字符串
	return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var claim = new(MyClaims)
	// 解析token 并将 payload 存放到 claim中
	token, err := jwt.ParseWithClaims(tokenString, claim, func(token *jwt.Token) (interface{}, error) {
		return mySecret, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, errors.New("token 签名无效")
		}
		return nil, errors.New("token 解析失败")
	}
	if !token.Valid {
		return nil, errors.New("token 无效")
	}
	return claim, nil
}
