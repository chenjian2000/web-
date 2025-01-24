package settings

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Conf 全库变量，用来保存程序的所有配置信息
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	*LogConfig   `mapstructure:"log"`
	*MySQLConfig `mapstructure:"mysql"`
	*RedisConfig `mapstructure:"redis"`
}

type LogConfig struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}

type MySQLConfig struct {
	Host         string `mapstructure:"host"`
	User         string `mapstructure:"user"`
	Port         int    `mapstructure:"port"`
	Password     string `mapstructure:"password"`
	DBname       string `mapstructure:"dbname"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
}

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	Port     int    `mapstructure:"port"`
	PoolSize int    `mapstructure:"pool_size"`
}

func Init() error {
	viper.SetConfigFile("./settings/config.yaml") // 使用正确相对路径
	err := viper.ReadInConfig()                   // 查找并读取配置文件
	if err != nil {
		// 处理读取配置文件的错误
		return fmt.Errorf("failed to read config: %w", err) // 包装错误信息
	}

	// 将配置信息反序列化到 Conf 结构体中
	if err := viper.Unmarshal(Conf); err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {
		fmt.Println("配置文件被修改！")
		if err := viper.Unmarshal(Conf); err != nil {
			fmt.Printf("failed to unmarshal config: %v\n", err)
		}
	})
	return nil
}
