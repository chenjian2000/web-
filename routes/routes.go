package routes

import (
	"net/http"
	"niko-web_app/controller"
	"niko-web_app/logger"
	"niko-web_app/settings"

	"github.com/gin-gonic/gin"
)

func Setup(cfg *settings.AppConfig) *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	// 注册业务 路由
	r.POST("/signup", controller.SignUpHandler)

	r.GET("/version", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to web_app",
			"version": cfg.Version,
		})
	})
	return r
}
