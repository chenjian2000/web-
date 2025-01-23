package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() error {
	viper.SetConfigFile("./settings/config.yaml") // 使用正确相对路径
	err := viper.ReadInConfig()                   // 查找并读取配置文件
	if err != nil {
		// 处理读取配置文件的错误
		return fmt.Errorf("failed to read config: %w", err) // 包装错误信息
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改！")
		err := viper.ReadInConfig() // 查找并读取配置文件
		if err != nil {
			// 处理读取配置文件的错误
			fmt.Printf("Failed to reload config: %v", err) // 移除错误返回
		}
	})
	return nil
}
