package mysql

import (
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"niko-web_app/models"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

const secret = "niko"

// CheckUserExist 根据用户名判断用户是否存在
func CheckUserExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExit
	}
	return
}

// InsertUser 像数据库中插入一条新的用户数据
func InsertUser(user *models.User) error {
	// 对密码进行加密
	password := encryptPassword(user.Password)
	// 执行 SQL 语句
	sqlStr := `insert into user(user_id, username, password) values(?,?,?)`
	_, err := db.Exec(sqlStr, user.UserID, user.Username, password)
	return err
}

// 数据加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

func Login(user *models.User) error {
	oPassword := user.Password
	sqlStr := `select user_id, username, password from user where username=?`
	err := db.Get(user, sqlStr, user.Username)
	if err == sql.ErrNoRows {
		return ErrorUserNotExit
	}
	if err != nil {
		// 查询数据库失败
		return err
	}
	password := encryptPassword(oPassword)
	if password != user.Password {
		return ErrorPasswordWrong
	}
	return nil
}
