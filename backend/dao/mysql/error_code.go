package mysql

import "errors"

// 定义业务状态
var (
	ErrorUserExit      = errors.New("用户已存在")
	ErrorUserNotExit   = errors.New("用户不存在")
	ErrorPasswordWrong = errors.New("密码错误")
	ErrorGenIDFailed   = errors.New("创建用户ID失败")
	ErrorInvalidID     = errors.New("无效的ID")
	ErrorQueryFailed   = errors.New("查询数据失败")
	ErrorInsertFailed  = errors.New("插入数据失败")
)
