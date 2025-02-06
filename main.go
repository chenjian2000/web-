package main

import (
	"context"
	"fmt"
	"net/http"
	"niko-web_app/dao/mysql"
	"niko-web_app/dao/redis"
	"niko-web_app/logger"
	"niko-web_app/pkg/snowflake"
	"niko-web_app/routes"
	"niko-web_app/settings"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"
)

// Go Web开发较通用的脚手架模板

func main() {
	// 1. 加载配置
	if err := settings.Init(); err != nil {
		fmt.Printf("init settings failed, err: %v\n", err)
		return
	}
	// 2. 初始化日志
	defer zap.L().Sync()
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("init logger failed, err: %v\n", err)
		return
	}
	// 3. 初始化MySQL连接
	if err := mysql.Init(settings.Conf.MySQLConfig); err != nil {
		fmt.Printf("init mysql failed, err: %v\n", err)
		return
	}
	defer mysql.Close()
	// 4. 初始化Redis连接
	if err := redis.Init(settings.Conf.RedisConfig); err != nil {
		fmt.Printf("init redis failed, err: %v\n", err)
		return
	}
	defer redis.Close()

	// 初始化 雪花算法分布式ID生成器
	if err := snowflake.Init(settings.Conf.StartTime, settings.Conf.MachineID); err != nil {
		fmt.Printf("init snowflake failed, err: %v\n", err)
		return
	}

	// 5. 注册路由
	router := routes.Setup(settings.Conf)
	// 6. 启动服务（优雅关机）
	// fmt.Println("启动服务器，监听端口：", settings.Conf.Port)
	// fmt.Println("version: ", settings.Conf.Version)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", settings.Conf.Port), // :8081
		Handler: router,
	} // 配置 http 服务器

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			zap.L().Fatal("listen: ", zap.Error(err))
		}
	}() // 启动服务器

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal, 1)
	// kill 默认会发送 syscall.SIGTERM 信号
	// kill -2 发送 syscall.SIGINT 信号，我们常用的Ctrl+C就是触发系统SIGINT信号
	// signal.Notify把收到的 syscall.SIGINT或syscall.SIGTERM 信号转发给quit
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) // 此处不会阻塞
	<-quit
	zap.L().Info("Shutdown Server ...")
	// 创建一个5秒超时的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// 5秒内优雅关闭服务（将未处理完的请求处理完再关闭服务），超过5秒就超时退出
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
