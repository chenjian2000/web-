package controller

type ResCode int64

const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParams
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeInvalidToken
	CodeInvalidAuthFormat
	CodeNotLogin
	ErrVoteRepeated
	ErrorVoteTimeExpire
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:         "success",
	CodeInvalidParams:   "请求参数错误",
	CodeUserExist:       "用户名重复",
	CodeUserNotExist:    "用户不存在",
	CodeInvalidPassword: "用户名或密码错误",
	CodeServerBusy:      "服务繁忙",

	CodeInvalidToken:      "无效的Token",
	CodeInvalidAuthFormat: "认证格式有误",
	CodeNotLogin:          "未登录",
	ErrVoteRepeated:       "请勿重复投票",
	ErrorVoteTimeExpire:   "投票时间已过",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if ok {
		return msg
	}
	return codeMsgMap[CodeServerBusy]
}
