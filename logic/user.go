package logic

import (
	"niko-web_app/dao/mysql"
	"niko-web_app/pkg/snowflake"
)

// 存放业务逻辑的处理

func SignUp() {
	// 判断用户是否存在
	mysql.QueryUserByName()
	// 生成 UID
	snowflake.GetID()
	// 写入数据库 mysql
	mysql.InsertUser()
}
