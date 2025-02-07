package controller

import (
	"net/http"
	"niko-web_app/logic"
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
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	// 2. 业务处理
	if err := logic.SignUp(&p); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"msg": "注册失败",
		})
		return
	}
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "success",
	})
}

// LoginHandler 处理用户登录
func LoginHandler(c *gin.Context) {
	// 1. 获取参数 和 校验参数
	var p models.ParamLogin
	if err := c.ShouldBindJSON(&p); err != nil {
		// 请求参数有误
		zap.L().Error("SignUp with invalid parameters", zap.Error(err))
		// 如果是 validator.ValidationErrors 类型的 错误
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusOK, gin.H{
				"msg": err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"msg": removeTopStruct(errs.Translate(trans)),
		})
		return
	}
	// 2. 业务处理
	if err := logic.Login(&p); err != nil {
		zap.L().Error("logic.Login failed", zap.Error(err))
		c.JSON(http.StatusOK, gin.H{
			"msg": "用户名或密码错误",
		})
		return
	}
	// 3. 返回响应
	c.JSON(http.StatusOK, gin.H{
		"msg": "登录成功",
	})
}
