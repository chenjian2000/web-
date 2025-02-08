package jwt

import "errors"

var (
	ErrorTokenSignatureInvalid = errors.New("token 签名无效")
	ErrorTokenParsingFailed    = errors.New("token 解析失败")
	ErrorTokenInvalid          = errors.New("token 无效")
)
