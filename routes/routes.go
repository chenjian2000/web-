package routes

import (
	"fmt"
	"net/http"
	"niko-web_app/controller"
	"niko-web_app/logger"
	"niko-web_app/pkg/jwt"
	"niko-web_app/settings"
	"strings"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *settings.AppConfig) *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务 路由
	r.POST("/signup", controller.SignUpHandler)

	// 登录业务 路由
	r.POST("/login", controller.LoginHandler)

	r.GET("/ping", JWTAuthMiddleware(), func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to web_app",
			"version": cfg.Version,
		})
	})
	return r
}

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 客户端携带Token有三种方式 1.放在请求头 2.放在请求体 3.放在URL
		// 这里假设Token放在Header的Authorization中，并使用Bearer开头: Authorization:Bearer XXX.XXXXX.XXXX
		// 这里的具体实现方式要依据你的实际业务情况决定
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, "请求头缺少Auth Token")
			c.Abort()
			return
		}
		// 按空格分割
		//&& parts[0] == "Bearer"
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2) {
			controller.ResponseErrorWithMsg(c, controller.CodeInvalidToken, "Token格式不对")
			c.Abort()
			return
		}
		// parts[1]是获取到的tokenString，我们使用之前定义好的解析JWT的函数来解析它
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			fmt.Println(err)
			controller.ResponseError(c, controller.CodeInvalidToken)
			c.Abort()
			return
		}
		// 将当前请求的userID信息保存到请求的上下文c上
		c.Set("userID", mc.UserID)
		c.Next() // 后续的处理函数可以用过c.Get(ContextUserIDKey)来获取当前请求的用户信息
	}
}
