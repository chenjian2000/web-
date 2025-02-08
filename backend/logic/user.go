package logic

import (
	"niko-web_app/dao/mysql"
	"niko-web_app/models"

	"niko-web_app/pkg/jwt"
	"niko-web_app/pkg/snowflake"
)

// 存放业务逻辑的处理

// SignUp 用户注册
func SignUp(p *models.ParamSignUp) (err error) {
	// 判断用户是否存在
	if err := mysql.CheckUserExist(p.Username); err != nil {
		return err
	}
	// 生成 UID
	userID, err := snowflake.GetID()
	if err != nil {
		return err
	}
	// 构造一个user实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	// 写入数据库 mysql
	return mysql.InsertUser(user)
}

// Login 用户登录
func Login(p *models.ParamLogin) (aToken, rToken string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	if err := mysql.Login(user); err != nil {
		return "", "", err
	}
	// 用户登录成功，生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
