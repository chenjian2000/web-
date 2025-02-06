package settings

import (
	"fmt"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	confOnce sync.Once
	initErr  error
)

// Conf 全局配置单例
var Conf = new(AppConfig)

type AppConfig struct {
	Name         string `mapstructure:"name"`
	Mode         string `mapstructure:"mode"`
	Version      string `mapstructure:"version"`
	StartTime    string `mapstructure:"start_time"`
	MachineID    uint16 `mapstructure:"machine_id"`
	Port         int    `mapstructure:"port"`
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
	confOnce.Do(func() {
		viper.SetConfigFile("./settings/config.yaml")
		if err := viper.ReadInConfig(); err != nil {
			initErr = fmt.Errorf("failed to read config: %w", err)
			return
		}

		if err := viper.Unmarshal(Conf); err != nil {
			initErr = fmt.Errorf("failed to unmarshal config: %w", err)
			return
		}

		viper.WatchConfig()
		viper.OnConfigChange(func(in fsnotify.Event) {
			fmt.Println("配置文件被修改！")
			if err := viper.Unmarshal(Conf); err != nil {
				fmt.Printf("failed to unmarshal config: %v\n", err)
			}
		})
	})
	return initErr
}
