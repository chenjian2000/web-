package controller

import (
	"errors"
	"net/http"
	"niko-web_app/dao/mysql"
	"niko-web_app/logic"
	"niko-web_app/pkg/jwt"
	"strings"

	"niko-web_app/models"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求的函数
func SignUpHandler(c *gin.Context) {
	// 1. 获取参数 和 参数校验
	var p models.ParamSignUp
	// 校验参数类型和格式
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误
		zap.L().Error("SignUp with invalid parameters", zap.Error(err))
		// 如果是 validator.ValidationErrors 类型的 错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(&p); err != nil {
		if errors.Is(err, mysql.ErrorUserExit) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3. 返回响应
	ResponseSuccess(c, nil)
}

// LoginHandler 处理用户登录
func LoginHandler(c *gin.Context) {
	// 1. 获取参数 和 校验参数
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误
		zap.L().Error("Login with invalid parameters", zap.Error(err))
		// 如果是 validator.ValidationErrors 类型的 错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, removeTopStruct(errs.Translate(trans)))
		return
	}
	// 2. 业务处理
	aToken, rToken, err := logic.Login(&p)
	if err != nil {
		zap.L().Error("logic.Login failed", zap.String("username", p.Username), zap.Error(err))
		if errors.Is(err, mysql.ErrorUserNotExit) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}

	// 3. 返回响应
	ResponseSuccess(c, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}

// RefreshTokenHandler 刷新 JWT
func RefreshTokenHandler(c *gin.Context) {
	rt := c.Query("refresh_token")
	// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URI
	// 这里假设Token放在Header的 Authorization 中，并使用 Bearer 开头
	// 这里的具体实现方式要依据你的实际业务情况决定
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		ResponseErrorWithMsg(c, CodeInvalidToken, "请求头缺少Auth Token")
		c.Abort()
		return
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		ResponseErrorWithMsg(c, CodeInvalidToken, "Token格式不对")
		c.Abort()
		return
	}
	aToken, rToken, err := jwt.RefreshToken(parts[1], rt)
	zap.L().Error("jwt.RefreshToken failed", zap.Error(err))
	c.JSON(http.StatusOK, gin.H{
		"access_token":  aToken,
		"refresh_token": rToken,
	})
}
