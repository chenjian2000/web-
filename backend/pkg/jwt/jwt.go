package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const AccessTokenExpireDuration = time.Hour * 2
const RefreshTokenExpireDuration = AccessTokenExpireDuration * 7

var mySecret = []byte("niko")

func keyFunc(token *jwt.Token) (interface{}, error) {
	return mySecret, nil
}

// 定义 payload （claims）
type MyClaims struct {
	UserID               uint64 `json:"user_id"`
	Username             string `json:"username"`
	jwt.RegisteredClaims        // jwt.StandardClaims 已经被弃用
}

// GenToken 生成 JWT token(access toekn) 和 refresh token
func GenToken(user_id uint64, username string) (aToken, rToken string, err error) {
	// 生成 payload 部分
	claims := &MyClaims{
		UserID:   user_id,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpireDuration)), // 过期时间
			Issuer:    "niko-goweb",                                                  // 签发人
		},
	}
	// 使用 HS256（HMAC SHA-256）签名算法创建 JWT 对象
	aToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(mySecret)

	// 生成refresh token
	rToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpireDuration)), // 过期时间
		Issuer:    "niko-goweb",                                                   // 签发人
	}).SignedString(mySecret)
	return
	// 使用密钥签名并生成最终的 JWT 字符串
	// return token.SignedString(mySecret)
}

func ParseToken(tokenString string) (*MyClaims, error) {
	var claim = new(MyClaims)
	// 解析token 并将 payload 存放到 claim中
	token, err := jwt.ParseWithClaims(tokenString, claim, keyFunc)
	if err != nil {
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return nil, ErrorTokenSignatureInvalid
		}
		return nil, ErrorTokenParsingFailed
	}
	if !token.Valid {
		return nil, ErrorTokenInvalid
	}
	return claim, nil
}

// RefreshToken 刷新AccessToken
func RefreshToken(aToken, rToken string) (newAToken, newRToken string, err error) {
	// refresh token无效直接返回
	if _, err = jwt.Parse(rToken, keyFunc); err != nil {
		return
	}

	// 从旧access token中解析出claims数据	解析出payload负载信息
	var claims MyClaims
	_, err = jwt.ParseWithClaims(aToken, &claims, keyFunc)
	v, _ := err.(*jwt.ValidationError)

	// 当access token是过期错误 并且 refresh token没有过期时就创建一个新的access token
	if v.Errors == jwt.ValidationErrorExpired {
		return GenToken(claims.UserID, claims.Username)
	}
	return
}
